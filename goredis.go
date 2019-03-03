package goredis // import "github.com/jbert/goredis"
import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
)

type Key string
type Val []byte

type ServerOp func(*Server, io.Writer) error

type Server struct {
	store Store
	ctx   context.Context
}

func (s *Server) Apply(conn net.Conn, op ServerOp) error {
	op(s, conn)
	return nil
}

func New() *Server {
	return &Server{
		store: NewMutexMapStore(),
		ctx:   context.Background(),
	}
}

func (s *Server) ListenAndServe(hostport string) error {
	log.Printf("Server starting - listening on %s", hostport)
	ln, err := net.Listen("tcp", hostport)
	if err != nil {
		return fmt.Errorf("Failed to listen: %s", err)
	}
	conn_id := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept: %s", err)
		}
		go s.handleConnection(context.WithValue(s.ctx, "id", conn_id), conn)
		conn_id += 1
	}

	return nil
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	id := ctx.Value("id").(int)
	log.Printf("CONN %d - open", id)
	defer func() {
		conn.Close()
		log.Printf("CONN %d - closed", id)
	}()
	op := func(s *Server, w io.Writer) error {
		err := s.store.Set(Key(fmt.Sprintf("key %d", id)), Val([]byte{1, 2, 3}))
		fmt.Fprintf(w, "Server has %d keys\n", s.store.Len())
		return err
	}
	s.Apply(conn, ServerOp(op))
}
