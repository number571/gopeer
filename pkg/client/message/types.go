package message

import "github.com/number571/go-peer/pkg/payload"

type IMessage interface {
	Head() iHead
	Body() iBody

	IsValid(uint64, uint64) bool
	Bytes() []byte
}

type iHead interface {
	Sender() []byte
	Session() []byte
	Salt() []byte
}

type iBody interface {
	Payload() payload.IPayload
	Hash() []byte
	Sign() []byte
	Proof() uint64
}
