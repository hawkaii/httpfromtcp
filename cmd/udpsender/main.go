package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	u, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, u)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	b := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("> ")
		line, err := b.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println(err)
			return
		}

		i, err := conn.Write([]byte(strings.TrimSpace(line)))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Sent %d bytes\n", i)

	}

}
