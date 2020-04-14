package gopeer

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"time"
)

// Return client hashname.
func (client *Client) Hashname() string {
	return client.hashname
}

// Return client's public key.
func (client *Client) Public() *rsa.PublicKey {
	x := *client.keys.public
	return &x
}

// Return client's private key.
func (client *Client) Private() *rsa.PrivateKey {
	x := *client.keys.private
	return &x
}

// Return listener address.
func (client *Client) Address() string {
	return client.address
}

// Return listener certificate.
func (client *Client) Certificate() []byte {
	return client.listener.certificate
}

// Return Destination struct from connected client.
func (client *Client) Destination(hash string) *Destination {
	if !client.InConnections(hash) {
		return nil
	}
	return &Destination{
		Address:     client.Connections[hash].address,
		Certificate: client.Connections[hash].certificate,
		Public:      client.Connections[hash].throwClient,
		Receiver:    client.Connections[hash].public,
	}
}

// Check if user saved in client data.
func (client *Client) InConnections(hash string) bool {
	if _, ok := client.Connections[hash]; ok {
		return true
	}
	return false
}

// Switcher function about GET and SET options.
// GET: accept package and send response;
// SET: accept package;
func (client *Client) HandleAction(title string, pack *Package, handleGet func(*Client, *Package) string, handleSet func(*Client, *Package)) bool {
	if pack.Head.Title != title {
		return false
	}
	switch pack.Head.Option {
	case settings.OPTION_GET:
		data := handleGet(client, pack)
		hash := pack.From.Sender.Hashname
		client.SendTo(client.Destination(hash), &Package{
			Head: Head{
				Title:  title,
				Option: settings.OPTION_SET,
			},
			Body: Body{
				Data: data,
			},
		})
	case settings.OPTION_SET:
		handleSet(client, pack)
	default:
		return false
	}
	return true
}

// Disconnect from user.
// Send package for disconnect.
// If the user is not responding: delete in local data.
func (client *Client) Disconnect(dest *Destination) error {
	var err error
	dest = client.wrapDest(dest)

	hash := HashPublic(dest.Receiver)
	if !client.InConnections(hash) {
		return errors.New("client not connected")
	}

	if client.Connections[hash].relation == nil {
		_, err = client.SendTo(dest, &Package{
			Head: Head{
				Title:  settings.TITLE_DISCONNECT,
				Option: settings.OPTION_GET,
			},
		})
	}

	if client.Connections[hash].relation != nil {
		client.Connections[hash].relation.Close()
	}

	delete(client.Connections, hash)
	return err
}

// Connect to user.
// Create local data with parameters.
// After sending GET and receiving SET package, set Connected = true.
func (client *Client) Connect(dest *Destination) error {
	dest = client.wrapDest(dest)
	var (
		session = GenerateRandomBytes(32)
		hash    = HashPublic(dest.Receiver)
	)
	if dest.Public == nil {
		return client.hiddenConnect(hash, session, dest.Receiver)
	}
	client.Connections[hash] = &Connect{
		connected:   false,
		hashname:    hash,
		address:     dest.Address,
		throwClient: dest.Public,
		public:      dest.Receiver,
		certificate: dest.Certificate,
		session:     session,
		Chans: Chans{
			Action: make(chan bool),
			action: make(chan bool),
		},
	}
	_, err := client.SendTo(dest, &Package{
		Head: Head{
			Title:  settings.TITLE_CONNECT,
			Option: settings.OPTION_GET,
		},
		Body: Body{
			Data: string(PackJSON(conndata{
				Certificate: Base64Encode(client.listener.certificate),
				Public:  Base64Encode([]byte(StringPublic(client.keys.public))),
				Session: Base64Encode(EncryptRSA(dest.Receiver, session)),
			})),
		},
	})
	if err != nil {
		return err
	}
	select {
	case <-client.Connections[hash].Chans.action:
		client.Connections[hash].connected = true
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		if client.Connections[hash].relation != nil {
			client.Connections[hash].relation.Close()
		}
		delete(client.Connections, hash)
		return errors.New("client not connected")
	}
	return nil
}

// Load file from node.
// Input = name file in node side.
// Output = result name file in our side.
func (client *Client) LoadFile(dest *Destination, input string, output string) error {
	dest = client.wrapDest(dest)

	hash := HashPublic(dest.Receiver)
	if !client.InConnections(hash) {
		return errors.New("client not connected")
	}

	if fileIsExist(output) {
		return errors.New("file already exists")
	}

	client.Connections[hash].transfer.active = true
	defer func() {
		client.Connections[hash].transfer.active = false
	}()

	for i := uint64(0) ;; i++ {
		client.SendTo(dest, &Package{
			Head: Head{
				Title:  settings.TITLE_FILETRANSFER,
				Option: settings.OPTION_GET,
			},
			Body: Body{
				Data: string(PackJSON(FileTransfer{
					Head: HeadTransfer{
						Id:   i,
						Name: input,
					},
				})),
			},
		})

		select {
		case <-client.Connections[hash].Chans.action:
			// pass
		case <-time.After(time.Duration(settings.WAITING_TIME * 2) * time.Second):
			return errors.New("waiting time is over")
		}

		var read = new(FileTransfer)
		UnpackJSON([]byte(client.Connections[hash].transfer.packdata), read)

		if read == nil {
			return errors.New("pack is null")
		}

		if read.Head.IsNull {
			break
		}

		data := read.Body.Data
		if !bytes.Equal(read.Body.Hash, HashSum(data)) {
			return errors.New("hash not equal file hash")
		}

		writeFile(output, read.Body.Data)
	}

	return nil
}

// Send by Destination.
func (client *Client) SendTo(dest *Destination, pack *Package) (*Package, error) {
	dest = client.wrapDest(dest)
	switch {
	case dest == nil: return nil, errors.New("dest is null")
	case dest.Public == nil: return nil, errors.New("public is null")
	case dest.Receiver == nil: return nil, errors.New("receiver is null")
	}

	pack.To.Receiver.Hashname = HashPublic(dest.Receiver)
	pack.To.Hashname = HashPublic(dest.Public)
	pack.To.Address = dest.Address

	return client.send(_confirm, pack)
}
