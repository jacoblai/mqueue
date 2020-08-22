package utils

import (
	"encoding/binary"
	"sync"
)

type Utils struct {
	pool *sync.Pool
}

func NewUtils() *Utils {
	return &Utils{
		pool: &sync.Pool{
			New: func() interface{} {
				cc := make([]byte, 8)
				return cc
			},
		},
	}
}

func (u *Utils) IdToKey(id uint64) []byte {
	key := u.pool.Get().([]byte)
	defer u.pool.Put(key)
	binary.BigEndian.PutUint64(key, id)
	return key
}

func (u *Utils) KeyToId(key []byte) uint64 {
	return binary.BigEndian.Uint64(key)
}
