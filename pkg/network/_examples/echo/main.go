package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

func main() {
	service := network.NewNode(nodeSettings(serviceAddress))
	service.HandleFunc(serviceHeader, func(_ network.INode, c conn.IConn, reqBytes []byte) {
		c.WritePayload(payload.NewPayload(
			serviceHeader,
			[]byte(fmt.Sprintf("echo: [%s]", string(reqBytes))),
		))
	})

	if err := service.Run(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // wait

	conn, err := conn.NewConn(
		connSettings(),
		serviceAddress,
	)
	if err != nil {
		panic(err)
	}

	pld, err := conn.FetchPayload(payload.NewPayload(
		serviceHeader,
		[]byte("hello, world!")),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pld.GetBody()))
}

func nodeSettings(serviceAddress string) network.ISettings {
	return network.NewSettings(&network.SSettings{
		FAddress:      serviceAddress,
		FCapacity:     (1 << 10),
		FMaxConnects:  1,
		FConnSettings: connSettings(),
	})
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FMessageSizeBytes: (1 << 10),
		FLimitVoidSize:    1, // not used
		FFetchTimeWait:    5 * time.Second,
	})
}
