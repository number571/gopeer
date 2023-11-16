package queue_set

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IQueueSet = &sQueueSet{}
)

type sQueueSet struct {
	fSettings ISettings
	fMutex    sync.Mutex
	fMap      map[string][]byte
	fQueue    []string
	fIndex    int
}

func NewQueueSet(pSettings ISettings) IQueueSet {
	return &sQueueSet{
		fSettings: pSettings,
		fQueue:    make([]string, pSettings.GetCapacity()),
		fMap:      make(map[string][]byte, pSettings.GetCapacity()),
	}
}

func (p *sQueueSet) GetSettings() ISettings {
	return p.fSettings
}

func (p *sQueueSet) Push(pKey, pValue []byte) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// hash already exists in queue
	sHash := encoding.HexEncode(pKey)
	if _, ok := p.fMap[sHash]; ok {
		return false
	}

	// delete old value in queue
	delete(p.fMap, p.fQueue[p.fIndex])

	// push hash to queue
	p.fQueue[p.fIndex] = sHash
	p.fMap[sHash] = pValue

	// increment queue index
	p.fIndex = (p.fIndex + 1) % len(p.fQueue)
	return true
}

func (p *sQueueSet) Load(pKey []byte) ([]byte, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	val, ok := p.fMap[encoding.HexEncode(pKey)]
	return val, ok
}
