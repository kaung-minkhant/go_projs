package main

import (
	"fmt"
	"time"
)

func main() {
  bigOperation()
}

func bigOperation() {
  defer trace("big operation")()
  defer trace("big")()
  time.Sleep(2 * time.Second)
}

func trace(name string) func() {
	fmt.Println("entering: ", name)
	return func() {
		fmt.Println("exiting: ", name)
	}
}
