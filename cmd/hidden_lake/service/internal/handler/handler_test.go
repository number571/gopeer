package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/closer"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/queue_set"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcServiceAddressInHLS = "hidden-echo-service"
	tcPathDBTemplate      = "database_test_%d.db"
	tcPathConfigTemplate  = "config_test_%d.yml"
)

var (
	tcConfig = fmt.Sprintf(`settings:
  message_size_bytes: 8192
  work_size_bits: 22
  key_size_bits: %d
  queue_period_ms: 1000
  limit_void_size_bytes: 4096
  network_key: test
address:
  tcp: test_address_tcp
  http: test_address_http
connections:
  - test_connect1
  - test_connect2
  - test_connect3
friends:
  test_recvr: %s
  test_name1: %s
  test_name2: %s
services:
  test_service1: 
    host: test_address1
  test_service2: 
    host: test_address2
  test_service3: 
    host: test_address3
`,
		testutils.TcKeySize,
		testutils.TgPubKeys[0],
		testutils.TgPubKeys[1],
		testutils.TgPubKeys[2],
	)
)

func testStartServerHTTP(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testEchoPage)

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: time.Second,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testEchoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		FMessage string `json:"message"`
	}

	var resp struct {
		FEcho  string `json:"echo"`
		FError int    `json:"error"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.FError = 1
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.FEcho = req.FMessage
	json.NewEncoder(w).Encode(resp)
}

func testAllCreate(cfgPath, dbPath, srvAddr string) (config.IWrapper, anonymity.INode, context.Context, context.CancelFunc, *http.Server) {
	wcfg := testNewWrapper(cfgPath)
	node, ctx, cancel := testRunNewNode(dbPath, "")
	srvc := testRunService(ctx, wcfg, node, srvAddr)
	time.Sleep(200 * time.Millisecond)
	return wcfg, node, ctx, cancel, srvc
}

func testAllFree(node anonymity.INode, cancel context.CancelFunc, srv *http.Server, pathCfg, pathDB string) {
	defer func() {
		os.RemoveAll(pathDB)
		os.RemoveAll(pathCfg)
	}()
	cancel()
	closer.CloseAll([]types.ICloser{
		srv,
		node.GetDBWrapper(),
		node.GetNetworkNode(),
	})
}

func testRunService(ctx context.Context, wcfg config.IWrapper, node anonymity.INode, addr string) *http.Server {
	mux := http.NewServeMux()

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	cfg := wcfg.GetConfig()
	mtx := &sync.Mutex{}

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleConfigSettingsPath, HandleConfigSettingsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, HandleConfigConnectsAPI(ctx, wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, HandleConfigFriendsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, HandleNetworkOnlineAPI(logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, HandleNetworkRequestAPI(ctx, mtx, cfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkPubKeyPath, HandleNetworkPubKeyAPI(logger, node))

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: time.Second,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testNewWrapper(cfgPath string) config.IWrapper {
	os.WriteFile(cfgPath, []byte(tcConfig), 0o644)
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	return config.NewWrapper(cfg)
}

func testRunNewNode(dbPath, addr string) (anonymity.INode, context.Context, context.CancelFunc) {
	os.RemoveAll(dbPath)
	node := testNewNode(dbPath, addr).HandleFunc(pkg_settings.CServiceMask, nil)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = node.Run(ctx) }()
	return node, ctx, cancel
}

func testNewNode(dbPath, addr string) anonymity.INode {
	db, err := database.NewKVDatabase(
		storage.NewSettings(&storage.SSettings{
			FPath:     dbPath,
			FWorkSize: testutils.TCWorkSize,
			FPassword: "CIPHER",
		}),
	)
	if err != nil {
		panic(err)
	}
	networkMask := uint64(1)
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   "TEST",
			FRetryEnqueue:  0,
			FNetworkMask:   networkMask,
			FFetchTimeWait: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		anonymity.NewDBWrapper().Set(db),
		testNewNetworkNode(addr),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: testutils.TCQueueCapacity,
				FVoidCapacity: testutils.TCQueueCapacity,
				FParallel:     1,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: testutils.TCMessageSize,
					FKeySizeBits:      testutils.TcKeySize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
			),
		).WithNetworkSettings(
			networkMask,
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
		),
		asymmetric.NewListPubKeys(),
	)
	return node
}

func testNewNetworkNode(addr string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FWorkSizeBits:     testutils.TCWorkSize,
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		}),
		queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: testutils.TCCapacity,
			}),
		),
	)
}
