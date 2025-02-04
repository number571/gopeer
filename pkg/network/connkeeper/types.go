package connkeeper

import (
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/types"
)

type IConnKeeper interface {
	types.IRunner

	GetSettings() ISettings
	GetNetworkNode() network.INode
}

type ISettings interface {
	GetConnections() []string
	GetDuration() time.Duration
}
