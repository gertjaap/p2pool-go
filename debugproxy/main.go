package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	listen := flag.Int("listen", 0, "The local port to listen on")
	remote := flag.String("remote", "", "The remote host to connect to")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *listen))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Waiting for client")

	client, err := l.Accept()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connecing to remote")

	conn, err := net.Dial("tcp", *remote)
	if err != nil {
		panic(err)
	}

	c2s, err := os.Create("client-to-server.log")
	defer c2s.Close()
	s2c, err := os.Create("server-to-client.log")
	defer s2c.Close()

	go func() {
		b := make([]byte, 1)
		for {
			i, err := client.Read(b)
			if err != nil {
				panic(err)
			}
			i2, err := conn.Write(b)
			if err != nil {
				panic(err)
			}
			if i != i2 {
				panic("Did not forward all bytes")
			}
			i3, err := c2s.Write([]byte(hex.EncodeToString(b)))
			if err != nil {
				panic(err)
			}
			if i2*2 != i3 {
				panic("Did not log all bytes")
			}
		}
	}()

	b := make([]byte, 1)
	for {
		i, err := conn.Read(b)
		if err != nil {
			panic(err)
		}
		i2, err := client.Write(b)
		if err != nil {
			panic(err)
		}
		if i != i2 {
			panic("Did not forward all bytes")
		}
		i3, err := s2c.Write([]byte(hex.EncodeToString(b)))
		if err != nil {
			panic(err)
		}
		if i2*2 != i3 {
			panic("Did not log all bytes")
		}
	}

}
