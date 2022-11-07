package network

import (
	"bytes"

	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/payload"
)

var (
	_ IMessage = sMessage{}
)

// Message = [Size[u64], Hash[uN], Head[u64], Body[u8...]]
type sMessage []byte

func NewMessage(pl payload.IPayload, key []byte) IMessage {
	payloadBytes := pl.Bytes()
	hashWithPayload := bytes.Join(
		[][]byte{
			hashing.NewHMACSHA256Hasher(key, payloadBytes).Bytes(),
			payloadBytes,
		},
		[]byte{},
	)

	return sMessage(bytes.Join(
		[][]byte{
			newPackage(hashWithPayload).SizeToBytes(),
			hashWithPayload,
		},
		[]byte{},
	))
}

func LoadMessage(bmsg []byte, key []byte) IMessage {
	if len(bmsg) < (cSizeUint + cSizeHash + cSizeHead) {
		return nil
	}

	// get Size[u64]
	mustLen := newPackage(bmsg[:cSizeUint]).BytesToSize()
	if mustLen != uint64(len(bmsg[cBeginHash:])) {
		return nil
	}

	// check Hash[uN]
	hashRecv := bmsg[cBeginHash:cEndHash]
	if !bytes.Equal(
		hashRecv,
		hashing.NewHMACSHA256Hasher(key, bmsg[cBeginHead:]).Bytes(),
	) {
		return nil
	}

	return sMessage(bmsg)
}

func (msg sMessage) Hash() []byte {
	return msg[cBeginHash:cEndHash]
}

func (msg sMessage) Payload() payload.IPayload {
	return payload.LoadPayload(msg[cBeginHead:])
}

func (msg sMessage) Bytes() []byte {
	return msg[:]
}
