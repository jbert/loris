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

type Sharded struct {
	shards []Store
}

func NewSharded(ctor func() Store) *Sharded {
	num_shards := 16
	ss := Sharded{
		shards: make([]Store, num_shards, num_shards),
	}
	for i := range ss.shards {
		ss.shards[i] = ctor()
	}
	return &ss
}

func (ss *Sharded) findShard(k Key) Store {
	// There are better hash functions
	var n rune
	for _, c := range k {
		n = n ^ c
		n = n >> 1
	}
	n = n % 16
	return ss.shards[n]
}

func (ss *Sharded) Set(k Key, v Val) error {
	return ss.findShard(k).Set(k, v)
}

func (ss *Sharded) Get(k Key) (Val, error) {
	return ss.findShard(k).Get(k)
}

func (ss *Sharded) Del(k Key) error {
	return ss.findShard(k).Del(k)
}

func (ss *Sharded) Keys() []Key {
	keys := make([]Key, 0)
	for _, s := range ss.shards {
		keys = append(keys, s.Keys()...)
	}
	return keys
}

func (ss *Sharded) Len() int {
	l := 0
	for i, s := range ss.shards {
		log.Printf("JB - shard %d has %d", i, s.Len())
		l += s.Len()
	}
	return l
}

type Map map[Key]Val

func NewMap() Map {
	return make(map[Key]Val)
}

func (ms Map) Set(k Key, v Val) error {
	m := map[Key]Val(ms)
	m[k] = v
	return nil
}

func (ms Map) Get(k Key) (Val, error) {
	m := map[Key]Val(ms)
	v, ok := m[k]
	if !ok {
		return nil, ErrNotExist
	}
	return v, nil
}

func (ms Map) Del(k Key) error {
	m := map[Key]Val(ms)
	_, exists := m[k]
	if !exists {
		return ErrNotExist
	}

	delete(ms, k)
	return nil
}

func (ms Map) Keys() []Key {
	m := map[Key]Val(ms)
	keys := make([]Key, 0, len(m))
	for k := range ms {
		keys = append(keys, k)
	}
	return keys
}

func (ms Map) Len() int {
	m := map[Key]Val(ms)
	return len(m)
}

type Mutex struct {
	sync.Mutex
	s Store
}

func NewMutex(s Store) *Mutex {
	return &Mutex{
		s: s,
	}
}

func (ms *Mutex) Set(k Key, v Val) error {
	ms.Lock()
	defer ms.Unlock()
	return ms.s.Set(k, v)
}

func (ms *Mutex) Get(k Key) (Val, error) {
	ms.Lock()
	defer ms.Unlock()
	return ms.s.Get(k)
}

func (ms *Mutex) Del(k Key) error {
	ms.Lock()
	defer ms.Unlock()
	return ms.s.Del(k)
}

func (ms *Mutex) Keys() []Key {
	ms.Lock()
	defer ms.Unlock()
	return ms.s.Keys()
}

func (ms *Mutex) Len() int {
	ms.Lock()
	defer ms.Unlock()
	return ms.s.Len()
}

type MutexMap struct {
	sync.Mutex
	m map[Key]Val
}

func NewMutexMap() *MutexMap {
	return &MutexMap{
		m: make(map[Key]Val),
	}
}

func (mms *MutexMap) Set(k Key, v Val) error {
	mms.Lock()
	defer mms.Unlock()

	mms.m[k] = v
	return nil
}

func (mms *MutexMap) Get(k Key) (Val, error) {
	mms.Lock()
	defer mms.Unlock()

	v, ok := mms.m[k]
	if !ok {
		return nil, ErrNotExist
	}
	return v, nil
}

func (mms *MutexMap) Del(k Key) error {
	mms.Lock()
	defer mms.Unlock()

	_, exists := mms.m[k]
	if !exists {
		return ErrNotExist
	}

	delete(mms.m, k)
	return nil
}

func (mms *MutexMap) Keys() []Key {
	mms.Lock()
	defer mms.Unlock()

	keys := make([]Key, 0, len(mms.m))
	for k := range mms.m {
		keys = append(keys, k)
	}
	return keys
}

func (mms *MutexMap) Len() int {
	mms.Lock()
	defer mms.Unlock()

	return len(mms.m)
}
