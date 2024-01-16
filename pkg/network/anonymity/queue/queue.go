package queue

import (
	"context"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IMessageQueue = &sMessageQueue{}
)

type sMessageQueue struct {
	fState state.IState
	fMutex sync.Mutex

	fSettings ISettings
	fClient   client.IClient

	fNetworkMask uint64
	fMsgSettings net_message.ISettings

	fMainPool *sMainPool
	fVoidPool *sVoidPool
}

type sMainPool struct {
	fMutex    sync.Mutex
	fQueue    chan net_message.IMessage
	fRawQueue chan message.IMessage
}

type sVoidPool struct {
	fMutex    sync.Mutex
	fQueue    chan net_message.IMessage
	fReceiver asymmetric.IPubKey
}

func NewMessageQueue(pSett ISettings, pClient client.IClient) IMessageQueue {
	return &sMessageQueue{
		fState:       state.NewBoolState(),
		fMsgSettings: net_message.NewSettings(&net_message.SSettings{}),
		fSettings:    pSett,
		fClient:      pClient,
		fMainPool: &sMainPool{
			fQueue:    make(chan net_message.IMessage, pSett.GetMainCapacity()),
			fRawQueue: make(chan message.IMessage, pSett.GetMainCapacity()),
		},
		fVoidPool: &sVoidPool{
			fQueue:    make(chan net_message.IMessage, pSett.GetVoidCapacity()),
			fReceiver: asymmetric.NewRSAPrivKey(pClient.GetPrivKey().GetSize()).GetPubKey(),
		},
	}
}

func (p *sMessageQueue) GetSettings() ISettings {
	return p.fSettings
}

func (p *sMessageQueue) GetClient() client.IClient {
	return p.fClient
}

func (p *sMessageQueue) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return utils.MergeErrors(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(nil) }()

	const numProcs = 2
	chBufErr := make(chan error, numProcs)

	wg := sync.WaitGroup{}
	wg.Add(numProcs)

	go p.runVoidPoolFiller(pCtx, &wg, chBufErr)
	go p.runMainPoolFiller(pCtx, &wg, chBufErr)

	wg.Wait()
	close(chBufErr)

	errList := make([]error, 0, numProcs)
	for err := range chBufErr {
		errList = append(errList, err)
	}
	return utils.MergeErrors(errList...)
}

func (p *sMessageQueue) runVoidPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
	defer pWg.Done()
	for {
		select {
		case <-pCtx.Done():
			chErr <- pCtx.Err()
			return
		default:
			if err := p.fillVoidPool(pCtx); err != nil {
				chErr <- err
				return
			}
		}
	}
}

func (p *sMessageQueue) runMainPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
	defer pWg.Done()
	for {
		select {
		case <-pCtx.Done():
			chErr <- pCtx.Err()
			return
		case x := <-p.fMainPool.fRawQueue:
			if err := p.fillMainPool(pCtx, x); err != nil {
				chErr <- err
				return
			}
		}
	}
}

func (p *sMessageQueue) WithNetworkSettings(pNetworkMask uint64, pMsgSettings net_message.ISettings) IMessageQueue {
	// stop all pools
	p.fMutex.Lock()
	p.fMainPool.fMutex.Lock()
	p.fVoidPool.fMutex.Lock()

	defer p.fMutex.Unlock()
	defer p.fMainPool.fMutex.Unlock()
	defer p.fVoidPool.fMutex.Unlock()

	// change net_message settings
	p.fNetworkMask = pNetworkMask
	p.fMsgSettings = pMsgSettings

	// clear all old queue state
	for len(p.fMainPool.fQueue) > 0 {
		<-p.fMainPool.fQueue
	}
	for len(p.fVoidPool.fQueue) > 0 {
		<-p.fVoidPool.fQueue
	}

	return p
}

func (p *sMessageQueue) EnqueueMessage(pMsg message.IMessage) error {
	p.fMainPool.fMutex.Lock()
	defer p.fMainPool.fMutex.Unlock()

	if p.mainPoolHasLimit() {
		return ErrQueueLimit
	}

	p.fMainPool.fRawQueue <- pMsg
	return nil
}

func (p *sMessageQueue) DequeueMessage(pCtx context.Context) net_message.IMessage {
	select {
	case <-pCtx.Done():
		return nil
	case <-time.After(p.fSettings.GetDuration()):
		select {
		case x := <-p.fMainPool.fQueue:
			// the main queue is checked first
			return x
		default:
			// take an existing message from any ready queue
			select {
			case <-pCtx.Done():
				return nil
			case x := <-p.fMainPool.fQueue:
				return x
			case x := <-p.fVoidPool.fQueue:
				return x
			}
		}
	}
}

func (p *sMessageQueue) fillMainPool(pCtx context.Context, pMsg message.IMessage) error {
	oldNetworkMask, oldMsgSettings := p.getNetworkSettings()
	chNetMsg := make(chan net_message.IMessage)
	go func() {
		chNetMsg <- net_message.NewMessage(
			oldMsgSettings,
			payload.NewPayload(
				oldNetworkMask,
				pMsg.ToBytes(),
			),
			p.fSettings.GetParallel(),
		)
	}()

	select {
	case <-pCtx.Done():
		return pCtx.Err()

	case netMsg := <-chNetMsg:
		newNetworkMask, newMsgSettings := p.getNetworkSettings()
		settingsChanged := newNetworkMask != oldNetworkMask ||
			newMsgSettings.GetNetworkKey() != oldMsgSettings.GetNetworkKey() ||
			newMsgSettings.GetWorkSizeBits() != oldMsgSettings.GetWorkSizeBits()

		if !settingsChanged {
			p.fMainPool.fQueue <- netMsg
		}
		return nil
	}
}

func (p *sMessageQueue) fillVoidPool(pCtx context.Context) error {
	if p.voidPoolHasLimit() {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fSettings.GetDuration() / 2):
			return nil
		}
	}

	msg, err := p.fClient.EncryptPayload(
		p.fVoidPool.fReceiver,
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		return err
	}

	oldNetworkMask, oldMsgSettings := p.getNetworkSettings()
	chNetMsg := make(chan net_message.IMessage)
	go func() {
		chNetMsg <- net_message.NewMessage(
			oldMsgSettings,
			payload.NewPayload(
				oldNetworkMask,
				msg.ToBytes(),
			),
			p.fSettings.GetParallel(),
		)
	}()

	select {
	case <-pCtx.Done():
		return pCtx.Err()

	case netMsg := <-chNetMsg:
		newNetworkMask, newMsgSettings := p.getNetworkSettings()
		settingsChanged := newNetworkMask != oldNetworkMask ||
			newMsgSettings.GetNetworkKey() != oldMsgSettings.GetNetworkKey() ||
			newMsgSettings.GetWorkSizeBits() != oldMsgSettings.GetWorkSizeBits()

		if !settingsChanged {
			p.fVoidPool.fQueue <- netMsg
		}
		return nil
	}
}

func (p *sMessageQueue) mainPoolHasLimit() bool {
	currLen1 := len(p.fMainPool.fQueue)
	currLen2 := len(p.fMainPool.fRawQueue)

	return false ||
		uint64(currLen1) >= p.fSettings.GetMainCapacity() ||
		uint64(currLen2) >= p.fSettings.GetMainCapacity()
}

func (p *sMessageQueue) voidPoolHasLimit() bool {
	p.fVoidPool.fMutex.Lock()
	defer p.fVoidPool.fMutex.Unlock()

	currLen := len(p.fVoidPool.fQueue)
	return uint64(currLen) >= p.fSettings.GetVoidCapacity()
}

func (p *sMessageQueue) getNetworkSettings() (uint64, net_message.ISettings) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fNetworkMask, p.fMsgSettings
}
