package database

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	gp_database "github.com/number571/go-peer/pkg/storage/database"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fDB    *gp_database.IKeyValueDB
}

func NewKeyValueDB(pPath string, pKey []byte) (IKeyValueDB, error) {
	db, err := gp_database.NewKeyValueDB(
		gp_database.NewSettings(&gp_database.SSettings{
			FPath:      pPath,
			FHashing:   true,
			FCipherKey: pKey,
		}),
	)
	if err != nil {
		return nil, errors.WrapError(err, "new key/value database")
	}
	return &sKeyValueDB{
		fDB: &db,
	}, nil
}

func (p *sKeyValueDB) Size(pR IRelation) uint64 {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.getSize(pR)
}

func (p *sKeyValueDB) Load(pR IRelation, pStart, pEnd uint64) ([]IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if pStart > pEnd {
		return nil, errors.NewError("start > end")
	}

	size := p.getSize(pR)
	if pEnd > size {
		return nil, errors.NewError("end > size")
	}

	res := make([]IMessage, 0, pEnd-pStart)
	for i := pStart; i < pEnd; i++ {
		data, err := (*p.fDB).Get(getKeyMessageByEnum(pR, i))
		if err != nil {
			return nil, errors.WrapError(err, "message undefined")
		}
		msg := LoadMessage(data)
		if msg == nil {
			return nil, errors.NewError("message is null")
		}
		res = append(res, msg)
	}

	return res, nil
}

func (p *sKeyValueDB) Push(pR IRelation, pMsg IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if _, err := (*p.fDB).Get(getKeyMessageByHash(pR, pMsg.GetSHA256UID())); err == nil {
		return errors.WrapError(err, "message is already exist")
	}

	if err := (*p.fDB).Set(getKeyMessageByHash(pR, pMsg.GetSHA256UID()), []byte{1}); err != nil {
		return errors.WrapError(err, "set uid to database")
	}

	size := p.getSize(pR)
	numBytes := encoding.Uint64ToBytes(size + 1)
	if err := (*p.fDB).Set(getKeySize(pR), numBytes[:]); err != nil {
		return errors.WrapError(err, "set size of message to database")
	}

	if err := (*p.fDB).Set(getKeyMessageByEnum(pR, size), pMsg.ToBytes()); err != nil {
		return errors.WrapError(err, "set message to database")
	}

	return nil
}

func (p *sKeyValueDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := (*p.fDB).Close(); err != nil {
		return errors.WrapError(err, "close KV database")
	}
	return nil
}

func (p *sKeyValueDB) getSize(pR IRelation) uint64 {
	data, err := (*p.fDB).Get(getKeySize(pR))
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}
