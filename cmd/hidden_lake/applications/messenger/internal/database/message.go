package database

import (
	"bytes"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	cIsIncomingSize = 1
	cTimestampSize  = encoding.CSizeUint64
	cPseudonymSize  = settings.CPseudonymSize
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fIsIncoming bool
	fTimestamp  uint64
	fPseudonym  []byte
	fMessage    []byte
}

func NewMessage(pIsIncoming bool, pPseudonym string, pMessage []byte) IMessage {
	if !utils.PseudonymIsValid(pPseudonym) {
		panic("!utils.PseudonymIsValid(pPseudonym)")
	}
	pseudonym := make([]byte, cPseudonymSize)
	copy(pseudonym, []byte(pPseudonym))
	return &sMessage{
		fIsIncoming: pIsIncoming,
		fTimestamp:  uint64(time.Now().Unix()),
		fPseudonym:  pseudonym,
		fMessage:    pMessage,
	}
}

func LoadMessage(pMsgBytes []byte) IMessage {
	if len(pMsgBytes) < (cIsIncomingSize + cTimestampSize + cPseudonymSize) {
		return nil
	}

	isIncoming := false
	if pMsgBytes[0] == 1 {
		isIncoming = true
	}

	blockTimestamp := [cTimestampSize]byte{}
	copy(blockTimestamp[:], pMsgBytes[cIsIncomingSize:cIsIncomingSize+cTimestampSize])

	return &sMessage{
		fIsIncoming: isIncoming,
		fTimestamp:  encoding.BytesToUint64(blockTimestamp),
		fPseudonym:  pMsgBytes[cIsIncomingSize+cTimestampSize : cIsIncomingSize+cTimestampSize+cPseudonymSize],
		fMessage:    pMsgBytes[cIsIncomingSize+cTimestampSize+cPseudonymSize:],
	}
}

func (p *sMessage) GetPseudonym() string {
	return string(bytes.TrimRight(p.fPseudonym, "\x00"))
}

func (p *sMessage) IsIncoming() bool {
	return p.fIsIncoming
}

func (p *sMessage) GetMessage() []byte {
	return p.fMessage
}

func (p *sMessage) GetTimestamp() string {
	return time.Unix(int64(p.fTimestamp), 0).Format("2006-01-02T15:04:05")
}

func (p *sMessage) ToBytes() []byte {
	firstByte := byte(0)
	if p.fIsIncoming {
		firstByte = byte(1)
	}
	blockTimestamp := encoding.Uint64ToBytes(p.fTimestamp)
	return bytes.Join(
		[][]byte{
			{firstByte},
			blockTimestamp[:],
			p.fPseudonym,
			p.fMessage,
		},
		[]byte{},
	)
}
