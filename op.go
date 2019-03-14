package loris

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jbert/loris/resp"
	"github.com/jbert/loris/store"
)

type OpFunc func(s store.Store, k store.Key, v store.Val) (store.Val, error)

func OpGet(s store.Store, k store.Key, v store.Val) (store.Val, error) {
	return s.Get(k)
}

func OpSet(s store.Store, k store.Key, v store.Val) (store.Val, error) {
	return v, s.Set(k, v)
}

func OpDel(s store.Store, k store.Key, v store.Val) (store.Val, error) {
	return nil, s.Del(k)
}

func ParseCommand(r io.Reader) (ServerOp, error) {
	respCmd, err := resp.Parse(bufio.NewScanner(r))
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
	bufs := make([][]byte, 0, len(a))
	for _, r := range a {
		buf, ok := r.([]byte)
		if !ok {
			return nil, fmt.Errorf("Command element is not []byte: %T", r)
		}
		bufs = append(bufs, buf)
	}

	switch strings.ToLower(string(bufs[0])) {
	case "command":
		op := func(s *Server, w io.Writer) (resp.Typ, bool) {
			return [][]byte{}, true
		}
		return op, nil
	case "del":
		op := func(s *Server, w io.Writer) (resp.Typ, bool) {
			num_deleted := 0
			for _, k := range bufs[1:] {
				err := s.store.Del(store.Key(k))
				if err == nil {
					num_deleted += 1
				}
			}
			return num_deleted, true
		}
		return op, nil
	case "dbsize":
		op := func(s *Server, w io.Writer) (resp.Typ, bool) {
			n := s.store.Len()
			return n, true
		}
		return op, nil
	case "get":
		op := func(s *Server, w io.Writer) (resp.Typ, bool) {
			v, err := s.store.Get(store.Key(bufs[1]))
			if err != nil {
				if err == store.ErrNotExist {
					return resp.TheNullBulkString, true
				}
				return err, true
			}
			return v, true
		}
		return op, nil
	case "set":
		op := func(s *Server, w io.Writer) (resp.Typ, bool) {
			err := s.store.Set(store.Key(bufs[1]), store.Val(bufs[2]))
			if err != nil {
				return err, true
			}
			//			return fmt.Sprintf("Server has %d keys", s.store.Len()), false
			return "OK", true
		}
		return op, nil
	case "quit":
		op := func(s *Server, w io.Writer) (resp.Typ, bool) {
			return "OK", false
		}
		return op, nil
	case "shutdown":
		op := func(s *Server, w io.Writer) (resp.Typ, bool) {
			s.StartShutdown()
			return "OK", false
		}
		return op, nil
	default:
		return nil, fmt.Errorf("Unrecognised command: %s", bufs[0])
	}
}
