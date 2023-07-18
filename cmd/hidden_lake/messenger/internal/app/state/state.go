package state

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/entropy"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

var (
	_ IStateManager = &sStateManager{}
	_ IClient       = &sClient{}
)

type sStateManager struct {
	fMutex    sync.Mutex
	fHashLP   []byte
	fConfig   config.IConfig
	fPathTo   string
	fStorage  storage.IKVStorage
	fDatabase database.IWrapperDB
	fClient   *sClient
}

type sClient struct {
	fService hls_client.IClient
	fTraffic hlt_client.IClient
}

func NewStateManager(
	pConfig config.IConfig,
	pPathTo string,
) IStateManager {
	stg, err := initCryptoStorage(pConfig, pPathTo)
	if err != nil {
		panic(err)
	}
	return &sStateManager{
		fConfig:   pConfig,
		fPathTo:   pPathTo,
		fStorage:  stg,
		fDatabase: database.NewWrapperDB(),
		fClient: &sClient{
			fService: hls_client.NewClient(
				hls_client.NewBuilder(),
				hls_client.NewRequester(
					fmt.Sprintf("http://%s", pConfig.GetConnection().GetService()),
					&http.Client{Timeout: time.Minute},
				),
			),
			fTraffic: hlt_client.NewClient(
				hlt_client.NewBuilder(),
				hlt_client.NewRequester(
					fmt.Sprintf("http://%s", pConfig.GetConnection().GetTraffic()),
					&http.Client{Timeout: time.Minute},
					message.NewSettings(&message.SSettings{
						FWorkSizeBits:     pConfig.GetWorkSizeBits(),
						FMessageSizeBytes: pConfig.GetMessageSizeBytes(),
					}),
				),
			),
		},
	}
}

func (p *sClient) Service() hls_client.IClient {
	return p.fService
}

func (p *sClient) Traffic() hlt_client.IClient {
	return p.fTraffic
}

func (p *sStateManager) GetConfig() config.IConfig {
	return p.fConfig
}

func (p *sStateManager) GetClient() IClient {
	return p.fClient
}

func (p *sStateManager) GetWrapperDB() database.IWrapperDB {
	return p.fDatabase
}

func (p *sStateManager) GetTemplate() *STemplateState {
	return &STemplateState{
		FAuthorized: p.StateIsActive(),
	}
}

func (p *sStateManager) CreateState(pHashLP []byte, pPrivKey asymmetric.IPrivKey) error {
	if _, err := p.fStorage.Get(pHashLP); err == nil {
		return errors.NewError("state already exists")
	}
	if err := p.newStorageState(pHashLP, pPrivKey); err != nil {
		return errors.WrapError(err, "new storage state")
	}
	return nil
}

func (p *sStateManager) OpenState(pHashLP []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.StateIsActive() {
		return errors.NewError("state already exists")
	}

	stateValue, err := p.getStorageState(pHashLP)
	if err != nil {
		return errors.WrapError(err, "get storage state")
	}

	entropyBooster := entropy.NewEntropyBooster(hlm_settings.CWorkForKeys, []byte{5, 7, 1})
	db, err := database.NewKeyValueDB(
		fmt.Sprintf("%s/%s", p.fPathTo, hlm_settings.CPathDB),
		entropyBooster.BoostEntropy(pHashLP),
	)
	if err != nil {
		return errors.WrapError(err, "open KV database")
	}

	p.GetWrapperDB().Set(db)
	if err := p.updateClientState(stateValue); err != nil {
		return errors.WrapError(err, "update client state")
	}
	p.fHashLP = pHashLP

	_ = p.updateClientTraffic(stateValue) // connect to HLT
	return nil
}

func (p *sStateManager) CloseState() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.StateIsActive() {
		return errors.NewError("state does not exist")
	}

	p.fHashLP = nil

	if err := p.GetWrapperDB().Close(); err != nil {
		return errors.WrapError(err, "close database")
	}

	if err := p.clearClientState(); err != nil {
		return errors.WrapError(err, "clear client state")
	}

	return nil
}

func (p *sStateManager) AddFriend(pAliasName string, pPubKey asymmetric.IPubKey) error {
	err := p.stateUpdater(
		p.updateClientFriends,
		func(storageValue *SStorageState) {
			storageValue.FFriends[pAliasName] = pPubKey.ToString()
		},
	)
	if err != nil {
		return errors.WrapError(err, "add friend (state updater)")
	}
	return nil
}

func (p *sStateManager) DelFriend(pAliasName string) error {
	err := p.stateUpdater(
		p.updateClientFriends,
		func(storageValue *SStorageState) {
			delete(storageValue.FFriends, pAliasName)
		},
	)
	if err != nil {
		return errors.WrapError(err, "del friend (state updater)")
	}
	return nil
}

func (p *sStateManager) AddConnection(pConnect string) error {
	err := p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = append(
				storageValue.FConnections,
				pConnect,
			)
		},
	)
	if err != nil {
		return errors.WrapError(err, "add connection (state updater)")
	}
	return nil
}

func (p *sStateManager) DelConnection(pConnect string) error {
	err := p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = remove(
				storageValue.FConnections,
				pConnect,
			)
		},
	)
	if err != nil {
		return errors.WrapError(err, "del connection (state updater)")
	}
	return nil
}

func (p *sStateManager) StateIsActive() bool {
	return p.fHashLP != nil
}

func (p *sStateManager) newStorageState(pHashLP []byte, pPrivKey asymmetric.IPrivKey) error {
	stateValueBytes := encoding.Serialize(&SStorageState{
		FPrivKey: pPrivKey.ToString(),
	}, false)
	if err := p.fStorage.Set(pHashLP, stateValueBytes); err != nil {
		return errors.WrapError(err, "new storage state")
	}
	return nil
}

func (p *sStateManager) setStorageState(pStateValue *SStorageState) error {
	stateValueBytes := encoding.Serialize(pStateValue, false)
	if err := p.fStorage.Set(p.fHashLP, stateValueBytes); err != nil {
		return errors.WrapError(err, "update storage state")
	}
	return nil
}

func (p *sStateManager) getStorageState(pHashLP []byte) (*SStorageState, error) {
	stateValueBytes, err := p.fStorage.Get(pHashLP)
	if err != nil {
		return nil, errors.WrapError(err, "get storage state bytes")
	}

	var stateValue = new(SStorageState)
	if err := encoding.Deserialize(stateValueBytes, stateValue); err != nil {
		return nil, errors.WrapError(err, "deserialize state")
	}

	return stateValue, nil
}

func remove(pSlice []string, pElem string) []string {
	for i, sElem := range pSlice {
		if pElem == sElem {
			return append(pSlice[:i], pSlice[i+1:]...)
		}
	}
	return pSlice
}
