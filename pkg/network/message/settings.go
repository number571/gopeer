package message

import "time"

var (
	_ IConstructSettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FTimestampWindow      time.Duration
	FWorkSizeBits         uint64
	FNetworkKey           string
	FParallel             uint64
	FRandMessageSizeBytes uint64
}

func NewSettings(pSett *SSettings) IConstructSettings {
	return (&sSettings{
		FTimestampWindow:      pSett.FTimestampWindow,
		FWorkSizeBits:         pSett.FWorkSizeBits,
		FNetworkKey:           pSett.FNetworkKey,
		FParallel:             pSett.FParallel,
		FRandMessageSizeBytes: pSett.FRandMessageSizeBytes,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() IConstructSettings {
	return p
}

func (p *sSettings) GetTimestampWindow() time.Duration {
	return p.FTimestampWindow
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *sSettings) GetParallel() uint64 {
	return p.FParallel
}

func (p *sSettings) GetRandMessageSizeBytes() uint64 {
	return p.FRandMessageSizeBytes
}
