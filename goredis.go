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

type ServerOp func(*Server, io.Writer) (error, bool)

type Server struct {
	store Store
	ctx   context.Context

	debug bool
}

func New() *Server {
	return &Server{
		store: NewMutexMapStore(),
		ctx:   context.Background(),
		debug: true,
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

type DebugReader struct {
	r io.Reader
}

func NewDebugReader(r io.Reader) *DebugReader {
	return &DebugReader{r: r}
}

func (dr *DebugReader) Read(p []byte) (int, error) {
	n, err := dr.r.Read(p)
	if err == nil {
		log.Printf("DR: [%d]: %02x [%s]", n, p[0:n], p[0:n])
	} else {
		log.Printf("DR: ERR [%s]", err)
	}
	return n, err
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	id := ctx.Value("id").(int)
	log.Printf("CONN %d - open", id)
	defer func() {
		conn.Close()
		log.Printf("CONN %d - closed", id)
	}()

	for {
		r := conn.(io.Reader)
		if s.debug {
			r = NewDebugReader(r)
		}
		op, err := ParseCommand(r)
		log.Printf("Op [%v] err [%v]", op, err)
		if err != nil {
			if err == io.EOF {
				break
			}
			// TODO: send resp
			fmt.Fprintf(conn, "-%s\r\n", err.Error())
			continue
		}

		err, ok := op(s, conn)
		if err != nil {
			// TODO: send resp
			fmt.Fprintf(conn, "-%s\r\n", err.Error())
		}
		if !ok {
			break
		}
	}
}
