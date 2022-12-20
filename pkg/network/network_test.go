package network

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcIter = 1000
)

func TestBroadcast(t *testing.T) {
	nodes, mapp := testNodes()
	defer testFreeNodes(nodes[:])

	// four receivers, sender not receive his messages
	tcMutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(4 * tcIter)

	headHandle := uint64(testutils.TcHead)
	handleF := func(node INode, conn conn.IConn, reqBytes []byte) {
		defer wg.Done()
		defer node.Broadcast(payload.NewPayload(headHandle, reqBytes))

		tcMutex.Lock()
		defer tcMutex.Unlock()

		val := string(reqBytes)
		flag, ok := mapp[node][val]
		if !ok {
			t.Errorf("incoming value '%s' undefined", val)
		}
		if flag {
			t.Errorf("incoming value '%s' already exists", val)
		}

		mapp[node][val] = true
	}

	for _, node := range nodes {
		node.Handle(headHandle, handleF)
	}

	// nodes[0] -> nodes[1:]
	for i := 0; i < tcIter; i++ {
		go func(i int) {
			pld := payload.NewPayload(
				headHandle,
				[]byte(fmt.Sprintf(testutils.TcLargeBodyTemplate, i)),
			)
			nodes[0].Broadcast(pld)
		}(i)
	}

	ch := make(chan struct{})
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()

	select {
	case <-ch:
	case <-time.After(20 * time.Second):
		t.Error("limit of waiting time for group")
		return
	}

	for _, node := range nodes {
		// pass sender
		if node == nodes[0] {
			continue
		}
		for i := 0; i < tcIter; i++ {
			val := fmt.Sprintf(testutils.TcLargeBodyTemplate, i)
			flag, ok := mapp[node][val]
			if !ok {
				t.Errorf("result value '%s' undefined", val)
				continue
			}
			if !flag {
				t.Errorf("result value '%s' not exists", val)
				continue
			}
		}
	}
}

func testNodes() ([5]INode, map[INode]map[string]bool) {
	nodes := [5]INode{}

	for i := 0; i < 5; i++ {
		sett := NewSettings(&SSettings{
			FCapacity:    (1 << 10),
			FMaxConnects: 10,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSize: (100 << 10),
				FTimeWait:    5 * time.Second,
			}),
		})
		nodes[i] = NewNode(sett)
	}

	go func() {
		err := nodes[2].Listen(testutils.TgAddrs[0])
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := nodes[4].Listen(testutils.TgAddrs[1])
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	nodes[0].Connect(testutils.TgAddrs[0])
	nodes[1].Connect(testutils.TgAddrs[1])

	nodes[3].Connect(testutils.TgAddrs[0])
	nodes[3].Connect(testutils.TgAddrs[1])

	mapp := make(map[INode]map[string]bool)
	for _, node := range nodes {
		// pass sender
		if node == nodes[0] {
			continue
		}
		mapp[node] = make(map[string]bool)
		for i := 0; i < tcIter; i++ {
			mapp[node][fmt.Sprintf(testutils.TcLargeBodyTemplate, i)] = false
		}
	}

	return nodes, mapp
}

func testFreeNodes(nodes []INode) {
	for _, node := range nodes {
		node.Close()
	}
}
