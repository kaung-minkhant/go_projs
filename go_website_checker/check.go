package main

import (
	"fmt"
	"net"
	"time"
)

func Check(domain, port string) string {
  address := fmt.Sprintf("%s:%s", domain, port)
  timeout := time.Duration(5 * time.Second)
  conn, err := net.DialTimeout("tcp", address, timeout)
  status := ""
  if err != nil {
    status = fmt.Sprintf("[DOWN] %v is unreachable, \n Error: %v\n", address, err)
  } else {
    status = fmt.Sprintf("[UP] %v is reachable, \n From: %v To: %v", address, conn.LocalAddr(), conn.RemoteAddr())
  }
  return status
}
