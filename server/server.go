package server

import (
	"log"
	"net"

	"github.com/bedakb/codecrafters-redis/handler"
	"github.com/bedakb/codecrafters-redis/store"
)

// Server holds the redis-server configuration.
type Server struct {
	storage *store.Store
	handler *handler.Handler
}

// ListenAndServe starts the server on the given port.
func ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	var srv Server
	srv.storage = store.New()
	srv.handler = handler.New(srv.storage)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go srv.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic:", r)
		}
	}()	
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("finished with err = %v", err)
			return
		}

		r, err := s.handler.Handle(buf[:n])
		if err != nil {
			log.Printf("finished with err = %v", err)
			return
		}

		conn.Write([]byte(r))
	}
}
