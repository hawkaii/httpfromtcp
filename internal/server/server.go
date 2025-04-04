package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/hawkaii/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener // The listener to accept connections
	closed   atomic.Bool  // Tracks if the server is closed
}

func Serve(port int) (*Server, error) {
	// Format the port into a string properly
	address := fmt.Sprintf(":%d", port)

	// Create the TCP listener
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	// Create the server struct and initialize it
	server := &Server{
		listener: l,
	}

	// Start accepting connections in a goroutine
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {

	err := s.listener.Close()

	s.closed.Store(true)

	return err

}

func (s *Server) listen() {
	for {

		conn, err := s.listener.Accept()

		if s.closed.Load() {
			return
		}
		if err != nil {
			fmt.Printf("nah")
			continue
		}

		go s.handle(conn)

	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	// Write the HTTP response
	err := response.WriteStatusLine(conn, response.StatusCode(200))

	if err != nil {
		return
	}

	h := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, h)
	if err != nil {
		return
	}

}
