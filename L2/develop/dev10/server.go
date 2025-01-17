package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func ServerStart() {
	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()

			r := bufio.NewReader(conn)
			w := bufio.NewWriter(conn)

			w.WriteString("Hi!\n")
			w.Flush()

			for {
				msg, err := r.ReadString('\n')
				if err != nil {
					break
				}

				msg = strings.TrimSpace(msg)
				fmt.Printf("New msg: %s\n", msg)

				w.WriteString(fmt.Sprintf("You send: %s\n", msg))
				w.Flush()

			}
		}(conn)
	}
}
