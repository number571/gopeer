package main

import (
	"fmt"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

func main() {
	var (
		client1 = newClient()
		client2 = newClient()
	)

	msg, err := client1.EncryptMessage(
		client2.GetPrivKey().GetPubKey(),
		payload.NewPayload64(0x0, []byte("hello, world!")).ToBytes(),
	)
	if err != nil {
		panic(err)
	}

	mapKeys := asymmetric.NewMapPubKeys(client1.GetPrivKey().GetPubKey())
	pubKey, decMsg, err := client2.DecryptMessage(mapKeys, msg)
	if err != nil {
		panic(err)
	}

	pld := payload.LoadPayload64(decMsg)
	fmt.Printf("Message: '%s';\nSender's public key: '%s';\n", string(pld.GetBody()), pubKey.ToString())
	fmt.Printf("Encrypted message: '%x'\n", msg)
}

func newClient() client.IClient {
	return client.NewClient(
		asymmetric.NewPrivKey(),
		(8 << 10),
	)
}
