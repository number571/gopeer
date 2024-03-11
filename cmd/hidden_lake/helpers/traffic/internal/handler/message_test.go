package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandleMessageAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[20]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, db)

	client := testNewClient()
	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg := net_message.NewMessage(
		testNetworkMessageSettings(),
		payload.NewPayload(hls_settings.CNetworkMask, msg.ToBytes()),
		1,
	)
	if err := hltClient.PutMessage(context.Background(), netMsg); err != nil {
		t.Error(err)
		return
	}

	strHash := encoding.HexEncode(netMsg.GetHash())
	gotNetMsg, err := hltClient.GetMessage(context.Background(), strHash)
	if err != nil {
		t.Error(err)
		return
	}

	gotEncMsg, err := message.LoadMessage(
		client.GetSettings(),
		gotNetMsg.GetPayload().GetBody(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	gotPubKey, gotPld, err := client.DecryptMessage(gotEncMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if gotPubKey.GetHasher().ToString() != client.GetPubKey().GetHasher().ToString() {
		t.Error(err)
		return
	}

	if string(gotPld.GetBody()) != testutils.TcBody {
		t.Error(err)
		return
	}
}
