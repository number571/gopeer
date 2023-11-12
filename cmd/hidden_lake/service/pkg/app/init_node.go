package app

import (
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/queue_pusher"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
)

func initNode(pCfg config.IConfig, pPrivKey asymmetric.IPrivKey, pLogger logger.ILogger) anonymity.INode {
	cfgSettings := pCfg.GetSettings()

	queueDuration := time.Duration(cfgSettings.GetQueuePeriodMS()) * time.Millisecond
	connDeadline := pkg_settings.GetConnDeadline(queueDuration)

	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   pkg_settings.CServiceName,
			FNetworkMask:   pkg_settings.CNetworkMask,
			FRetryEnqueue:  pkg_settings.CRetryEnqueue,
			FFetchTimeWait: pkg_settings.CFetchTimeout,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		pLogger,
		anonymity.NewWrapperDB(),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      pCfg.GetAddress().GetTCP(),
				FMaxConnects:  pkg_settings.CNetworkMaxConns,
				FReadTimeout:  queueDuration,
				FWriteTimeout: queueDuration,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:       pCfg.GetNetworkKey(),
					FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
					FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
					FLimitVoidSize:    cfgSettings.GetLimitVoidSizeBytes(),
					FWaitReadDeadline: pkg_settings.CConnWaitReadDeadline,
					FReadDeadline:     connDeadline,
					FWriteDeadline:    connDeadline,
				}),
			}),
			queue_pusher.NewQueuePusher(
				queue_pusher.NewSettings(&queue_pusher.SSettings{
					FCapacity: pkg_settings.CNetworkQueueSize,
				}),
			),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: pkg_settings.CQueueCapacity,
				FPoolCapacity: pkg_settings.CQueuePoolCapacity,
				FDuration:     queueDuration,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
					FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
				}),
				pPrivKey,
			),
		),
		func() asymmetric.IListPubKeys {
			f2f := asymmetric.NewListPubKeys()
			for _, pubKey := range pCfg.GetFriends() {
				f2f.AddPubKey(pubKey)
			}
			return f2f
		}(),
	)
}
