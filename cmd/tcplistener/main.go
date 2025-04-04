package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkaii/httpfromtcp/internal/request"
)

func main() {

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	// Notify the channel for SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer l.Close()

	// Start a goroutine to handle the signal
	go func() {
		<-sigChan
		fmt.Println("\nShutting down server...")
		l.Close()
		os.Exit(0)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Accepted connection")

		req, err := request.RequestFromReader(conn)

		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
		}

		fmt.Println("Request line: ")

		fmt.Printf(
			"- Method: %s\n- Target: %s\n- Version: %s\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)

		fmt.Println("Headers:")

		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Println("Body: ")

		fmt.Printf("%s\n", string(req.Body))

		conn.Close()

	}

}
