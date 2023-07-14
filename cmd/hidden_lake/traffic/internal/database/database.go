package database

import (
	"sync"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
)

type sKeyValueDB struct {
	fMutex   sync.Mutex
	fPointer uint64

	fSettings ISettings
	fDB       database.IKVDatabase
}

func NewKeyValueDB(pSett ISettings) (IKVDatabase, error) {
	if pSett.GetCapacity() == 0 {
		return nil, errors.NewError("capacity of messages = 0")
	}

	kvDB, err := database.NewKeyValueDB(
		storage.NewSettings(&storage.SSettings{
			FPath:      pSett.GetPath(),
			FHashing:   false,
			FCipherKey: []byte("_"),
		}),
	)
	if err != nil {
		return nil, errors.WrapError(err, "new key/value database")
	}

	db := &sKeyValueDB{
		fSettings: pSett,
		fDB:       kvDB,
	}
	db.fPointer = db.getPointer()
	return db, nil
}

func (p *sKeyValueDB) GetOriginal() database.IKVDatabase {
	return p.fDB
}

func (p *sKeyValueDB) Settings() ISettings {
	return p.fSettings
}

func (p *sKeyValueDB) Hashes() ([]string, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	msgsLimit := p.Settings().GetCapacity()
	res := make([]string, 0, msgsLimit)
	for i := uint64(0); i < msgsLimit; i++ {
		hash, err := p.fDB.Get(getKeyHash(i))
		if err != nil {
			break
		}
		if len(hash) != hashing.CSHA256Size {
			return nil, errors.NewError("incorrect hash size")
		}
		res = append(res, encoding.HexEncode(hash))
	}

	return res, nil
}

func (p *sKeyValueDB) Push(pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	hash := pMsg.GetBody().GetHash()
	if _, err := p.fDB.Get(getKeyMessage(hash)); err == nil {
		return nil
	}

	params := message.NewSettings(&message.SSettings{
		FWorkSizeBits:     p.Settings().GetWorkSizeBits(),
		FMessageSizeBytes: p.Settings().GetMessageSizeBytes(),
	})
	if !pMsg.IsValid(params) {
		return errors.NewError("invalid push message")
	}

	// delete old message
	keyHash := getKeyHash(p.getPointer())
	if hash, err := p.fDB.Get(keyHash); err == nil {
		keyMsg := getKeyMessage(hash)
		if err := p.fDB.Del(keyMsg); err != nil {
			return errors.WrapError(err, "delete old key")
		}
	}

	// rewrite hash's field
	newHash := pMsg.GetBody().GetHash()
	if err := p.fDB.Set(keyHash, newHash); err != nil {
		return errors.WrapError(err, "rewrite key hash")
	}

	// write message
	keyMsg := getKeyMessage(newHash)
	if err := p.fDB.Set(keyMsg, pMsg.ToBytes()); err != nil {
		return errors.WrapError(err, "write message")
	}

	// update pointer
	if err := p.incPointer(); err != nil {
		return errors.WrapError(err, "increment pointer")
	}

	return nil
}

func (p *sKeyValueDB) Load(pStrHash string) (message.IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	hash := encoding.HexDecode(pStrHash)
	if len(hash) != hashing.CSHA256Size {
		return nil, errors.NewError("key size invalid")
	}

	data, err := p.fDB.Get(getKeyMessage(hash))
	if err != nil {
		return nil, errors.NewError("message undefined")
	}

	msg := message.LoadMessage(
		message.NewSettings(&message.SSettings{
			FWorkSizeBits:     p.Settings().GetWorkSizeBits(),
			FMessageSizeBytes: p.Settings().GetMessageSizeBytes(),
		}),
		data,
	)
	if msg == nil {
		panic("message is nil")
	}

	return msg, nil
}

func (p *sKeyValueDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return errors.WrapError(err, "close KV database")
	}
	return nil
}

func (p *sKeyValueDB) getPointer() uint64 {
	data, err := p.fDB.Get(getKeyPointer())
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}

func (p *sKeyValueDB) incPointer() error {
	msgsLimit := p.Settings().GetCapacity()
	res := encoding.Uint64ToBytes((p.getPointer() + 1) % msgsLimit)
	if err := p.fDB.Set(getKeyPointer(), res[:]); err != nil {
		return errors.WrapError(err, "set pointer into KV database")
	}
	return nil
}
