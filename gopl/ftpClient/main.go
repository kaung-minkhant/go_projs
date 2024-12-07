package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strconv"
)

func usage() {
	fmt.Println(`./ftpClient <port> <command>`)
}

var allowedCommands = []string{
	"ls",
}

func main() {
	if len(os.Args) < 3 {
		usage()
		log.Fatal("bad arguments")
	}

	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		usage()
		log.Fatal("invalid port")
	}

	if !slices.Contains(allowedCommands, os.Args[2]) {
		usage()
		log.Fatal("invalid command")
	}

	conn, err := net.Dial("tcp", "localhost:"+os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	command := os.Args[2]

	switch command {
	case "ls":
		_, err = conn.Write([]byte("ls"))
		if err != nil {
			log.Fatal(err)
		}
	}

}
