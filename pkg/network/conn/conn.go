package conn

import (
	"bytes"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	// IV + Hash + PayloadHead
	cPayloadSizeOverHead = symmetric.CAESBlockSize + hashing.CSHA256Size + encoding.CSizeUint64
)

const (
	cWorkSize = 1
	cSalt     = "_"
)

var (
	_ IConn = &sConn{}
)

type sConn struct {
	fMutex      sync.Mutex
	fSocket     net.Conn
	fSettings   ISettings
	fCipher     symmetric.ICipher
	fNetworkKey string
}

func NewConn(pSett ISettings, pAddr string) (IConn, error) {
	conn, err := net.Dial("tcp", pAddr)
	if err != nil {
		return nil, errors.WrapError(err, "tcp connect")
	}
	return LoadConn(pSett, conn), nil
}

func LoadConn(pSett ISettings, pConn net.Conn) IConn {
	keyBuilder := keybuilder.NewKeyBuilder(cWorkSize, []byte(cSalt))
	cipherKey := keyBuilder.Build([]byte(pSett.GetNetworkKey()))
	return &sConn{
		fSettings:   pSett,
		fSocket:     pConn,
		fNetworkKey: pSett.GetNetworkKey(),
		fCipher:     symmetric.NewAESCipher(cipherKey),
	}
}

func (p *sConn) GetNetworkKey() string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fNetworkKey
}

func (p *sConn) SetNetworkKey(pNetworkKey string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	keyBuilder := keybuilder.NewKeyBuilder(cWorkSize, []byte(cSalt))
	cipherKey := keyBuilder.Build([]byte(pNetworkKey))
	p.fCipher = symmetric.NewAESCipher(cipherKey)
}

func (p *sConn) GetSettings() ISettings {
	return p.fSettings
}

func (p *sConn) GetSocket() net.Conn {
	return p.fSocket
}

func (p *sConn) FetchPayload(pPld payload.IPayload) (payload.IPayload, error) {
	if err := p.WritePayload(pPld); err != nil {
		return nil, errors.WrapError(err, "write payload")
	}

	chPld := make(chan payload.IPayload)
	go p.readPayload(chPld)

	select {
	case rpld := <-chPld:
		if rpld == nil {
			return nil, errors.NewError("read payload")
		}
		return rpld, nil
	case <-time.After(p.fSettings.GetFetchTimeWait()):
		return nil, errors.NewError("read payload (timeout)")
	}
}

func (p *sConn) Close() error {
	return p.fSocket.Close()
}

func (p *sConn) WritePayload(pPld payload.IPayload) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	encMsgBytes := p.fCipher.EncryptBytes(
		message.NewMessage(
			pPld,
			[]byte(p.fNetworkKey),
		).ToBytes(),
	)

	prng := random.NewStdPRNG()
	voidBytes := prng.GetBytes(prng.GetUint64() % (p.fSettings.GetLimitVoidSize() + 1))

	err := p.sendBytes(bytes.Join(
		[][]byte{
			// send headers with length of blocks
			p.getBlockSize(encMsgBytes),
			p.getBlockSize(voidBytes),
			// send data blocks
			encMsgBytes,
			voidBytes,
		},
		[]byte{},
	))
	if err != nil {
		return errors.WrapError(err, "send payload bytes")
	}

	return nil
}

func (p *sConn) ReadPayload() payload.IPayload {
	chPld := make(chan payload.IPayload)
	go p.readPayload(chPld)
	return <-chPld
}

func (p *sConn) sendBytes(pBytes []byte) error {
	bytesPtr := uint64(len(pBytes))
	for {
		p.fSocket.SetWriteDeadline(time.Now().Add(p.fSettings.GetWriteDeadline()))
		n, err := p.fSocket.Write(pBytes[:bytesPtr])
		if err != nil {
			return errors.WrapError(err, "write tcp bytes")
		}

		bytesPtr = bytesPtr - uint64(n)
		pBytes = pBytes[:bytesPtr]

		if bytesPtr == 0 {
			break
		}
	}
	return nil
}

func (p *sConn) getBlockSize(pBytes []byte) []byte {
	blockSize := encoding.Uint64ToBytes(uint64(len(pBytes)))
	return p.fCipher.EncryptBytes(blockSize[:])
}

func (p *sConn) recvBlockSize(deadline time.Duration) (uint64, error) {
	p.fSocket.SetReadDeadline(time.Now().Add(deadline))

	encBufLen := make([]byte, symmetric.CAESBlockSize+encoding.CSizeUint64)
	n, err := p.fSocket.Read(encBufLen)
	if err != nil {
		return 0, errors.WrapError(err, "read tcp block size")
	}

	if n != symmetric.CAESBlockSize+encoding.CSizeUint64 {
		return 0, errors.NewError("block size is invalid")
	}

	// mustLen = Size[u64] in uint64
	bufLen := p.fCipher.DecryptBytes(encBufLen)
	arrLen := [encoding.CSizeUint64]byte{}
	copy(arrLen[:], bufLen)

	return encoding.BytesToUint64(arrLen), nil
}

func (p *sConn) readPayload(pChPld chan payload.IPayload) {
	var pld payload.IPayload
	defer func() {
		pChPld <- pld
	}()

	msgSize, err := p.recvBlockSize(p.fSettings.GetWaitReadDeadline()) // the connection has not sent anything yet
	if err != nil || msgSize > (p.fSettings.GetMessageSizeBytes()+cPayloadSizeOverHead) {
		return
	}

	voidSize, err := p.recvBlockSize(p.fSettings.GetReadDeadline())
	if err != nil || voidSize > p.fSettings.GetLimitVoidSize() {
		return
	}

	mustLen := msgSize + voidSize
	dataRaw := make([]byte, 0, mustLen)
	for {
		p.fSocket.SetReadDeadline(time.Now().Add(p.fSettings.GetReadDeadline()))

		buffer := make([]byte, mustLen)
		n, err := p.fSocket.Read(buffer)
		if err != nil {
			return
		}

		dataRaw = bytes.Join(
			[][]byte{
				dataRaw,
				buffer[:n],
			},
			[]byte{},
		)

		mustLen -= uint64(n)
		if mustLen == 0 {
			break
		}
	}

	// try unpack message from bytes
	msg := message.LoadMessage(
		p.getCipher().DecryptBytes(dataRaw[:msgSize]),
		[]byte(p.getNetworkKey()),
	)
	if msg == nil {
		return
	}

	pld = msg.GetPayload()
}

func (p *sConn) getNetworkKey() string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fNetworkKey
}

func (p *sConn) getCipher() symmetric.ICipher {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fCipher
}
