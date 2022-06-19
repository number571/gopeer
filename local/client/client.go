package client

import (
	"bytes"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/puzzle"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/crypto/symmetric"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/routing"
	"github.com/number571/go-peer/settings"
)

var (
	_ IClient = &sClient{}
)

// Basic structure describing the user.
type sClient struct {
	fSettings settings.ISettings
	fPrivKey  asymmetric.IPrivKey
}

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(priv asymmetric.IPrivKey, sett settings.ISettings) IClient {
	if priv == nil {
		return nil
	}
	return &sClient{
		fSettings: sett,
		fPrivKey:  priv,
	}
}

// Get public key from client object.
func (client *sClient) PubKey() asymmetric.IPubKey {
	return client.PrivKey().PubKey()
}

// Get private key from client object.
func (client *sClient) PrivKey() asymmetric.IPrivKey {
	return client.fPrivKey
}

// Get settings from client object.
func (client *sClient) Settings() settings.ISettings {
	return client.fSettings
}

// Function wrap message in multiple route encryption.
// Need use pseudo sender if route not null.
func (client *sClient) Encrypt(route routing.IRoute, msg message.IMessage) (message.IMessage, []byte) {
	var (
		psender       = NewClient(route.PSender(), client.Settings())
		rmsg, session = client.onceEncrypt(route.Receiver(), msg)
	)
	if psender == nil && len(route.List()) != 0 {
		return nil, nil
	}
	for _, pub := range route.List() {
		rmsg, _ = psender.(*sClient).onceEncrypt(
			pub,
			message.NewMessage(
				encoding.Uint64ToBytes(client.Settings().Get(settings.CMaskRout)),
				rmsg.ToPackage().Bytes(),
			),
		)
	}
	return rmsg, session
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *sClient) onceEncrypt(receiver asymmetric.IPubKey, msg message.IMessage) (message.IMessage, []byte) {
	var (
		rand    = random.NewStdPRNG()
		salt    = rand.Bytes(client.Settings().Get(settings.CSizeSkey))
		session = rand.Bytes(client.Settings().Get(settings.CSizeSkey))
	)

	data := bytes.Join(
		[][]byte{
			encoding.Uint64ToBytes(uint64(len(msg.Body().Data()))),
			msg.Body().Data(),
			encoding.Uint64ToBytes(rand.Uint64() % (settings.CSizePack / 4)),
		},
		[]byte{},
	)

	hash := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			salt,
			client.PubKey().Bytes(),
			receiver.Bytes(),
			data,
		},
		[]byte{},
	)).Bytes()

	cipher := symmetric.NewAESCipher(session)
	return &message.SMessage{
		FHead: message.SHeadMessage{
			FSender:  cipher.Encrypt(client.PubKey().Bytes()),
			FSession: receiver.Encrypt(session),
			FSalt:    cipher.Encrypt(salt),
		},
		FBody: message.SBodyMessage{
			FData:  cipher.Encrypt(data),
			FHash:  hash,
			FSign:  cipher.Encrypt(client.PrivKey().Sign(hash)),
			FProof: puzzle.NewPoWPuzzle(client.Settings().Get(settings.CSizeWork)).Proof(hash),
		},
	}, session
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *sClient) Decrypt(msg message.IMessage) (message.IMessage, []byte) {
	const (
		SizeUint64 = 8 // bytes
	)

	if msg == nil {
		return nil, nil
	}

	// Initial check.
	if len(msg.Body().Hash()) != hashing.GSHA256Size {
		return nil, nil
	}

	// Proof of work. Prevent spam.
	diff := client.Settings().Get(settings.CSizeWork)
	puzzle := puzzle.NewPoWPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil, nil
	}

	// Decrypt session key by private key of receiver.
	session := client.PrivKey().Decrypt(msg.Head().Session())
	if session == nil {
		return nil, nil
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := symmetric.NewAESCipher(session)
	publicBytes := cipher.Decrypt(msg.Head().Sender())
	if publicBytes == nil {
		return nil, nil
	}

	// Load public key and check standart size.
	public := asymmetric.LoadRSAPubKey(publicBytes)
	if public == nil {
		return nil, nil
	}
	if public.Size() != client.PubKey().Size() {
		return nil, nil
	}

	// Decrypt salt.
	salt := cipher.Decrypt(msg.Head().Salt())
	if salt == nil {
		return nil, nil
	}

	// Decrypt main data of message by session key.
	dataBytes := cipher.Decrypt(msg.Body().Data())
	if dataBytes == nil {
		return nil, nil
	}

	// Check received hash and generated hash.
	check := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			salt,
			publicBytes,
			client.PubKey().Bytes(),
			dataBytes,
		},
		[]byte{},
	)).Bytes()
	if !bytes.Equal(check, msg.Body().Hash()) {
		return nil, nil
	}

	// check size of data
	if len(dataBytes) < SizeUint64 {
		return nil, nil
	}

	// pass random bytes and get main data
	mustLen := encoding.BytesToUint64(dataBytes[:SizeUint64])
	allData := dataBytes[SizeUint64:]
	if mustLen > uint64(len(allData)) {
		return nil, nil
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.Decrypt(msg.Body().Sign())
	if sign == nil {
		return nil, nil
	}
	if !public.Verify(msg.Body().Hash(), sign) {
		return nil, nil
	}

	decMsg := &message.SMessage{
		FHead: message.SHeadMessage{
			FSender:  publicBytes,
			FSession: session,
			FSalt:    salt,
		},
		FBody: message.SBodyMessage{
			FData:  allData[:mustLen],
			FHash:  msg.Body().Hash(),
			FSign:  sign,
			FProof: msg.Body().Proof(),
		},
	}

	// export title from (title||data)
	title := exportMessage(decMsg)
	if title == nil {
		return nil, nil
	}

	// Return decrypted message with title.
	return decMsg, title
}

// used in decrypt function from client
// export title from (title||data)
// store (title||data) <- data
func exportMessage(msg *message.SMessage) []byte {
	const (
		SizeUint64 = 8 // bytes
	)

	if len(msg.FBody.FData) < SizeUint64 {
		return nil
	}

	mustLen := encoding.BytesToUint64(msg.FBody.FData[:SizeUint64])
	allData := msg.FBody.FData[SizeUint64:]
	if mustLen > uint64(len(allData)) {
		return nil
	}

	msg.FBody.FData = allData[mustLen:]
	return allData[:mustLen]
}
