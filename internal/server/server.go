package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/hawkaii/httpfromtcp/internal/request"
	"github.com/hawkaii/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(io.Writer, *request.Request) *HandlerError

func writeHandlerError(s response.StatusCode, msg string) *HandlerError {
	h := &HandlerError{
		StatusCode: s,
		Message:    msg,
	}
	return h
}

type Server struct {
	listener net.Listener // The listener to accept connections
	closed   atomic.Bool  // Tracks if the server is closed
}

func Serve(port int, handler Handler) (*Server, error) {
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
	go server.listen(handler)

	return server, nil
}

func (s *Server) Close() error {

	err := s.listener.Close()

	s.closed.Store(true)

	return err

}

func (s *Server) listen(handler Handler) {
	for {

		conn, err := s.listener.Accept()

		if s.closed.Load() {
			return
		}
		if err != nil {
			fmt.Printf("nah")
			continue
		}

		go s.handle(conn, handler)

	}

}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()
	// Write the HTTP response

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	buf := bytes.Buffer{}

	handlerErr := handler(&buf, req)
	if handlerErr != nil {
		err = response.WriteStatusLine(conn, handlerErr.StatusCode)
		fmt.Printf("%d \n message: %s", handlerErr.StatusCode, handlerErr.Message)
		if err != nil {
			fmt.Println(err)
			return
		}

		errorHeaders := response.GetDefaultHeaders(len(handlerErr.Message))
		err = response.WriteHeaders(conn, errorHeaders)
		if err != nil {
			fmt.Printf("Error writing header: %s\n", err)
			return
		}

		_, err = conn.Write([]byte(handlerErr.Message))
		if err != nil {
			fmt.Printf("Error writing header: %s\n", err)
			return
		}
	} else {
		err = response.WriteStatusLine(conn, response.StatusCodeSuccess)
		if err != nil {
			fmt.Printf("Error writing status line: %s\n", err)
			return
		}

		h := response.GetDefaultHeaders(buf.Len())

		err = response.WriteHeaders(conn, h)
		if err != nil {
			fmt.Printf("Error writing header: %s\n", err)
			return
		}

		_, err = io.Copy(conn, &buf)
		if err != nil {
			fmt.Printf("Error writing header: %s\n", err)
			return
		}

	}

}
