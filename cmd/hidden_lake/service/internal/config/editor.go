package config

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
)

var (
	_ IEditor = &sEditor{}
)

type sEditor struct {
	fMutex  sync.Mutex
	fConfig *SConfig
}

func newEditor(pCfg IConfig) IEditor {
	if pCfg == nil {
		return nil
	}
	v, ok := pCfg.(*SConfig)
	if !ok {
		return nil
	}
	return &sEditor{
		fConfig: v,
	}
}

func (p *sEditor) UpdateNetworkKey(pNetworkKey string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.WrapError(err, "load config (update connections)")
	}

	cfg := icfg.(*SConfig)
	cfg.FNetworkKey = pNetworkKey
	err = filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg, true))
	if err != nil {
		return errors.WrapError(err, "write config (update connections)")
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FNetworkKey = cfg.FNetworkKey
	return nil
}

func (p *sEditor) UpdateConnections(pConns []string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.WrapError(err, "load config (update connections)")
	}

	cfg := icfg.(*SConfig)
	cfg.FConnections = deleteDuplicateStrings(pConns)
	err = filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg, true))
	if err != nil {
		return errors.WrapError(err, "write config (update connections)")
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FConnections = cfg.FConnections
	return nil
}

func (p *sEditor) UpdateFriends(pFriends map[string]asymmetric.IPubKey) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	for name, pubKey := range pFriends {
		if pubKey.GetSize() == p.fConfig.GetKeySizeBits() {
			continue
		}
		return errors.NewError(fmt.Sprintf("not supported key size for '%s'", name))
	}

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.WrapError(err, "load config (update friends)")
	}

	cfg := icfg.(*SConfig)
	cfg.fFriends = deleteDuplicatePubKeys(pFriends)
	cfg.FFriends = pubKeysToStrings(pFriends)
	err = filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg, true))
	if err != nil {
		return errors.WrapError(err, "write config (update friends)")
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.fFriends = cfg.fFriends
	p.fConfig.FFriends = cfg.FFriends
	return nil
}

func pubKeysToStrings(pPubKeys map[string]asymmetric.IPubKey) map[string]string {
	result := make(map[string]string, len(pPubKeys))
	for name, pubKey := range pPubKeys {
		result[name] = pubKey.ToString()
	}
	return result
}

func deleteDuplicatePubKeys(pPubKeys map[string]asymmetric.IPubKey) map[string]asymmetric.IPubKey {
	result := make(map[string]asymmetric.IPubKey, len(pPubKeys))
	mapping := make(map[string]struct{})
	for name, pubKey := range pPubKeys {
		pubStr := pubKey.GetAddress().ToString()
		if _, ok := mapping[pubStr]; ok {
			continue
		}
		mapping[pubStr] = struct{}{}
		result[name] = pubKey
	}
	return result
}

func deleteDuplicateStrings(pStrs []string) []string {
	result := make([]string, 0, len(pStrs))
	mapping := make(map[string]struct{})
	for _, s := range pStrs {
		if _, ok := mapping[s]; ok {
			continue
		}
		mapping[s] = struct{}{}
		result = append(result, s)
	}
	return result
}
