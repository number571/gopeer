package main

import (
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
)

const (
	networkKey  = "network-key"
	networkMask = uint64(0x1122334455667788)
	keySize     = uint64(1024)
	msgSize     = uint64(8192)
	workSize    = uint64(10)
)

func newNode(serviceName, address string) anonymity.INode {
	db, err := database.NewKVDatabase(
		storage.NewSettings(&storage.SSettings{
			FPath: "./database_" + serviceName + ".db",
		}),
	)
	if err != nil {
		panic(err)
	}

	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   serviceName,
			FNetworkMask:   networkMask,
			FFetchTimeWait: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{
				FInfo: os.Stdout,
				FWarn: os.Stdout,
				FErro: os.Stdout,
			}),
			func(ia logger.ILogArg) string {
				logGetterFactory, ok := ia.(anon_logger.ILogGetterFactory)
				if !ok {
					panic("got invalid log arg")
				}

				logGetter := logGetterFactory.Get()
				if logGetter.GetHash() == nil {
					// request with null hash from sender
					return ""
				}

				return fmt.Sprintf(
					"name=%s code=%02x hash=%X proof=%08d bytes=%d",
					logGetter.GetService(),
					logGetter.GetType(),
					logGetter.GetHash()[:16],
					logGetter.GetProof(),
					logGetter.GetSize(),
				)
			},
		),
		anonymity.NewDBWrapper().Set(db),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      address,
				FMaxConnects:  256,
				FReadTimeout:  time.Minute,
				FWriteTimeout: time.Minute,
				FConnSettings: newConnSettings(msgSize, workSize, networkKey),
			}),
			lru.NewLRUCache(
				lru.NewSettings(&lru.SSettings{
					FCapacity: 1024,
				}),
			),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FDuration:     2 * time.Second,
				FParallel:     1,
				FMainCapacity: 32,
				FVoidCapacity: 32,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      keySize,
					FMessageSizeBytes: msgSize,
				}),
				asymmetric.NewRSAPrivKey(keySize),
			),
		).WithNetworkSettings(
			networkMask,
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: workSize,
				FNetworkKey:   networkKey,
			}),
		),
		asymmetric.NewListPubKeys(),
	)
}

func newConnSettings(mSize, wSize uint64, nKey string) conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FNetworkKey:       nKey,
		FWorkSizeBits:     wSize,
		FMessageSizeBytes: mSize,
		FLimitVoidSize:    8192,
		FWaitReadDeadline: time.Hour,
		FReadDeadline:     time.Minute,
		FWriteDeadline:    time.Minute,
	})
}
