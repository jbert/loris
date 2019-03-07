package store

import (
	"errors"
	"log"
	"sync"
)

type Key string
type Val []byte

var ErrNotExist = errors.New("Key does not exist")

type Store interface {
	Set(k Key, v Val) error
	Get(k Key) (Val, error)
	Del(k Key) error
	// TODO: keys isn't performance critical but this interface sucks
	Keys() []Key
	Len() int
}

type ShardedStore struct {
	shards []Store
}

func NewShardedStore(ctor func() Store) *ShardedStore {
	num_shards := 16
	ss := ShardedStore{
		shards: make([]Store, num_shards, num_shards),
	}
	for i := range ss.shards {
		ss.shards[i] = ctor()
	}
	return &ss
}

func (ss *ShardedStore) findShard(k Key) Store {
	// There are better hash functions
	var n rune
	for _, c := range k {
		n = n ^ c
		n = n >> 1
	}
	n = n % 16
	return ss.shards[n]
}

func (ss *ShardedStore) Set(k Key, v Val) error {
	return ss.findShard(k).Set(k, v)
}

func (ss *ShardedStore) Get(k Key) (Val, error) {
	return ss.findShard(k).Get(k)
}

func (ss *ShardedStore) Del(k Key) error {
	return ss.findShard(k).Del(k)
}

func (ss *ShardedStore) Keys() []Key {
	keys := make([]Key, 0)
	for _, s := range ss.shards {
		keys = append(keys, s.Keys()...)
	}
	return keys
}

func (ss *ShardedStore) Len() int {
	l := 0
	for i, s := range ss.shards {
		log.Printf("JB - shard %d has %d", i, s.Len())
		l += s.Len()
	}
	return l
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

	_, exists := mms.m[k]
	if !exists {
		return ErrNotExist
	}

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

func (mms *MutexMapStore) Len() int {
	mms.Lock()
	defer mms.Unlock()

	return len(mms.m)
}
