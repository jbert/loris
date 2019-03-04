package goredis

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jbert/goredis/resp"
)

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

func ParseCommand(r io.Reader) (ServerOp, error) {
	respCmd, err := resp.Parse(bufio.NewReader(r))
	if err != nil {
		return nil, err
	}
	return decodeRespCommand(respCmd)
}

func decodeRespCommand(rCmd resp.Typ) (ServerOp, error) {
	a, ok := rCmd.([]resp.Typ)
	if !ok {
		return nil, fmt.Errorf("Command is not array: %T", rCmd)
	}
	ss := make([]string, 0, len(a))
	for _, r := range a {
		s, ok := r.(string)
		if !ok {
			return nil, fmt.Errorf("Command element is not string: %T", s)
		}
		ss = append(ss, s)
	}

	switch strings.ToLower(ss[0]) {
	case "set":
		op := func(s *Server, w io.Writer) (error, bool) {
			err := s.store.Set(Key(ss[1]), Val(ss[2]))
			fmt.Fprintf(w, "Server has %d keys\n", s.store.Len())
			return err, false
		}
		return op, nil
	default:
		return nil, fmt.Errorf("Unrecnogised command: %s", ss[0])
	}
}
