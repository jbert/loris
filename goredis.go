package goredis // import "github.com/jbert/goredis"
import (
	"fmt"
	"log"
	"net"
)

type Key string
type Val []byte

type Server struct {
	store Store
}

func New() *Server {
	return &Server{
		store: NewMutexMapStore(),
	}
}

func (s *Server) ListenAndServe(hostport string) error {
	ln, err := net.Listen("tcp", hostport)
	if err != nil {
		return fmt.Errorf("Failed to listen: %s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept: %s", err)
		}
		go s.handleConnection(conn)
	}

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Printf("New connection")
	defer conn.Close()
}
