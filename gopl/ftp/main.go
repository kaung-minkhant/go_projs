package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func usage() {
	fmt.Println(`./ftp <port>`)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("bad arguments")
	}
	listner, err := net.Listen("tcp", "localhost:"+os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Printf("error occured when establishing connection: %s\n", err)
			continue
		}

		go handleConn(conn)
	}
}

func mustCopy(dest io.Writer, src io.Reader) {
	if _, err := io.Copy(dest, src); err != nil {
		log.Fatal(err)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	var buf bytes.Buffer
	mustCopy(&buf, conn)
	fmt.Printf("data from %s: %s\n", conn.RemoteAddr(), buf.String())
}
