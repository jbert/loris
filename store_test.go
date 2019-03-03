package goredis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TK1 = Key("foo")
var TK2 = Key("why hello there: 💩")
var TV1 = Val([]byte{0, 1, 2, 3})
var TV2 = Val([]byte{4, 5, 6, 7, 8, 9, 10})

func TestMutexMapStore(t *testing.T) {
	mms := NewMutexMapStore()
	t.Run("MutexMapStore", func(t *testing.T) { testAll(t, mms) })
}

func testAll(t *testing.T, s Store) {
	t.Run("basic", func(t *testing.T) { testBasic(t, s) })
}

func testBasic(t *testing.T, s Store) {
	a := assert.New(t)
	a.Equal([]Key{}, s.Keys())

	a.NoError(s.Set(TK1, TV1))
	v, err := s.Get(TK1)
	a.NoError(err)
	a.Equal(TV1, v)
	a.Equal([]Key{TK1}, s.Keys())

	a.NoError(s.Set(TK2, TV2))
	v, err = s.Get(TK2)
	a.NoError(err)
	a.Equal(TV2, v)
	a.Equal([]Key{TK1, TK2}, s.Keys())

	a.NoError(s.Del(TK1))
	v, err = s.Get(TK1)
	a.EqualError(err, ErrNotExist.Error(), ErrNotExist)
	a.Equal([]Key{TK2}, s.Keys())
}
