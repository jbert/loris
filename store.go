package goredis

import (
	"errors"
	"sync"
)

var ErrNotExist = errors.New("Key does not exist")

type Store interface {
	Set(k Key, v Val) error
	Get(k Key) (Val, error)
	Del(k Key) error
	// TODO: keys isn't performance critical but this interface sucks
	Keys() []Key
}

type MutexMapStore struct {
	sync.Mutex
	m map[Key]Val
}

func NewMutexMapStore() *MutexMapStore {
	return &MutexMapStore{
		m: make(map[Key]Val),
	}
}

func (mms *MutexMapStore) Set(k Key, v Val) error {
	mms.Lock()
	defer mms.Unlock()

	mms.m[k] = v
	return nil
}

func (mms *MutexMapStore) Get(k Key) (Val, error) {
	mms.Lock()
	defer mms.Unlock()

	v, ok := mms.m[k]
	if !ok {
		return nil, ErrNotExist
	}
	return v, nil
}

func (mms *MutexMapStore) Del(k Key) error {
	mms.Lock()
	defer mms.Unlock()

	delete(mms.m, k)
	return nil
}

func (mms *MutexMapStore) Keys() []Key {
	mms.Lock()
	defer mms.Unlock()

	keys := make([]Key, 0, len(mms.m))
	for k := range mms.m {
		keys = append(keys, k)
	}
	return keys
}
