package goredis

type OpFunc func(s Store, k Key, v Val) (Val, error)

func OpGet(s Store, k Key, v Val) (Val, error) {
	return s.Get(k)
}

func OpSet(s Store, k Key, v Val) (Val, error) {
	return v, s.Set(k, v)
}

func OpDel(s Store, k Key, v Val) (Val, error) {
	return nil, s.Del(k)
}
