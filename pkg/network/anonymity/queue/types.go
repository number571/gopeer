package queue

import (
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
)

type IMessageQueue interface {
	types.IApp
	WithNetworkSettings(uint64, net_message.ISettings) IMessageQueue

	GetSettings() ISettings
	GetClient() client.IClient

	EnqueueMessage(message.IMessage) error
	DequeueMessage() net_message.IMessage
}

type ISettings interface {
	GetMainCapacity() uint64
	GetPoolCapacity() uint64
	GetDuration() time.Duration
}
