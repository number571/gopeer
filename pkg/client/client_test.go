package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcMessageSize = (2 << 10)
)

var (
	tgPrivKey  = asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024)
	tgMsgLimit = GetMessageLimit(
		tcMessageSize,       // 10KiB
		tgPrivKey.GetSize(), // 1024bits
	)
	tgMessages = []string{
		testutils.TcBody,
		"",
		"A",
		"AA",
		"AAA",
		"AAAA",
		"AAAAA",
		"AAAAAA",
		"AAAAAAA",
		"AAAAAAAA",
		"AAAAAAAAA",
		"AAAAAAAAAA",
		"AAAAAAAAAAA",
		"AAAAAAAAAAAA",
		"AAAAAAAAAAAAA",
		"AAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		random.NewStdPRNG().GetString(tgMsgLimit), // maximum size of message
	}
)

func TestEncrypt(t *testing.T) {
	client1 := testNewClient()
	client2 := testNewClient()

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg, err := client1.EncryptPayload(client2.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}

	_, decPl, err := client2.DecryptMessage(msg)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal([]byte(testutils.TcBody), decPl.GetBody()) {
		t.Error("data not equal with decrypted data")
		return
	}
}

func TestMessageSize(t *testing.T) {
	client1 := testNewClient()
	sizes := make([]int, 0, len(tgMessages))

	for _, smsg := range tgMessages {
		pl := payload.NewPayload(uint64(testutils.TcHead), []byte(smsg))
		msg, err := client1.EncryptPayload(client1.GetPubKey(), pl)
		if err != nil {
			t.Error(err)
			return
		}
		sizes = append(sizes, len(msg.ToBytes()))
	}

	for i := 0; i < len(sizes)-1; i++ {
		if sizes[i] != sizes[i+1] {
			t.Errorf(
				"len bytes of different messages = id(%d, %d) not equals = size(%d, %d)",
				i, i+1,
				sizes[i], sizes[i+1],
			)
			return
		}
	}
}

func TestGetMessageLimit(t *testing.T) {
	client1 := testNewClient()

	msg1 := random.NewStdPRNG().GetBytes(tgMsgLimit)
	pld1 := payload.NewPayload(uint64(testutils.TcHead), []byte(msg1))
	if _, err := client1.EncryptPayload(client1.GetPubKey(), pld1); err != nil {
		t.Error("message1 > message limit")
		return
	}

	msg2 := random.NewStdPRNG().GetBytes(tgMsgLimit + 1)
	pld2 := payload.NewPayload(uint64(testutils.TcHead), []byte(msg2))
	if _, err := client1.EncryptPayload(client1.GetPubKey(), pld2); err == nil {
		t.Error("message2 > message limit but not alert")
		return
	}
}

func testNewClient() IClient {
	return NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSizeBits:     testutils.TCWorkSize,
			FMessageSizeBytes: tcMessageSize,
		}),
		tgPrivKey,
	)
}
