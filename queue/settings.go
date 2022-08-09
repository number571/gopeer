package queue

import "time"

type sSettings struct {
	fMainCapacity uint64
	fPullCapacity uint64
	fMessageSize  uint64
	fDuration     time.Duration
}

func NewSettings(mCapacity, pCapacity, msgSize uint64, duration time.Duration) ISettings {
	return &sSettings{
		fMainCapacity: mCapacity,
		fPullCapacity: pCapacity,
		fDuration:     duration,
		fMessageSize:  msgSize,
	}
}

func (s *sSettings) GetMainCapacity() uint64 {
	return s.fMainCapacity
}

func (s *sSettings) GetPullCapacity() uint64 {
	return s.fPullCapacity
}

func (s *sSettings) GetMessageSize() uint64 {
	return s.fMessageSize
}

func (s *sSettings) GetDuration() time.Duration {
	return s.fDuration
}
