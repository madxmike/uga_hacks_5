package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {

	go StartServer("3000")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
