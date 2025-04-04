package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkaii/httpfromtcp/internal/request"
	"github.com/hawkaii/httpfromtcp/internal/response"
	"github.com/hawkaii/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(w io.Writer, req *request.Request) *server.HandlerError {

		if req.RequestLine.RequestTarget == "/yourproblem" {

			return &server.HandlerError{
				StatusCode: response.StatusCodeBadRequest,
				Message:    "Your problem is not my problem\n",
			}

		}

		if req.RequestLine.RequestTarget == "/myproblem" {

			return &server.HandlerError{
				StatusCode: response.StatusCodeInternalServerError,
				Message:    "Woopsie, my bad\n",
			}

		}

		_, err := w.Write([]byte("All good, frfr\n"))
		if err != nil {
			return &server.HandlerError{
				StatusCode: response.StatusCodeInternalServerError,
				Message:    "Error writing response\n",
			}
		}

		return nil

	}
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
