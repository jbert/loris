package loris // import "github.com/jbert/loris"
import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jbert/loris/resp"
	"github.com/jbert/loris/store"
)

type ServerOp func(*Server, io.Writer) (resp.Typ, bool)

type Server struct {
	store  store.Store
	ctx    context.Context
	cancel context.CancelFunc

	debug bool
}

func NewWithStore(s store.Store, debug bool) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		store:  s,
		debug:  debug,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) ListenAndServe(hostport string) error {
	log.Printf("Server starting - listening on %s", hostport)
	ln, err := net.Listen("tcp", hostport)
	if err != nil {
		return fmt.Errorf("Failed to listen: %s", err)
	}
	connID := 0
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Failed to accept: %s", err)
			}
			go s.handleConnection(context.WithValue(s.ctx, "id", connID), conn)
			connID ++
		}
	}()
	<-s.ctx.Done()

	return nil
}

func (s *Server) StartShutdown() {
	s.cancel()
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
	if s.debug {
		log.Printf("CONN %d - open", id)
	}
	defer func() {
		conn.Close()
		if s.debug {
			log.Printf("CONN %d - closed", id)
		}
	}()

	for {
		r := conn.(io.Reader)
		if s.debug {
			r = NewDebugReader(r)
		}
		op, err := ParseCommand(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			conn.Write(resp.TypToBuf(fmt.Errorf("Parsing error: %s", err)))
			continue
		}

		rsp, ok := op(s, conn)
		conn.Write(resp.TypToBuf(rsp))
		if !ok {
			break
		}
	}
}
