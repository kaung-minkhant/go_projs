package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	clocks := make(map[string]string)
	for _, clock := range os.Args[1:] {
		parts := strings.Split(clock, "=")
		if len(parts) != 2 {
			log.Fatal("bad arguments")
		}
		clocks[parts[0]] = parts[1]
	}

	for name, host := range clocks {
		go watch(name, host)
	}

	for {
		time.Sleep(1 * time.Second)
	}

	// conns := make([]net.Conn, 0)
	// for _, clock := range clocks {
	//   conn, err := net.Dial("tcp", clock)
	//   if err != nil {
	//     log.Fatal(err)
	//   }
	//   conns = append(conns, conn)
	// }

}

func watch(name, host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return
	}
	defer conn.Close()

	r := bufio.NewScanner(conn)
	for r.Scan() {
		fmt.Printf("%s: %s\n", name, r.Text())
	}
}

func mustCopy(dest io.Writer, src io.Reader) {
	if _, err := io.Copy(dest, src); err != nil {
		log.Fatal(err)
	}
}
