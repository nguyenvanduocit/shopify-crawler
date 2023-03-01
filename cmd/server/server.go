package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server, cleanup, err := InitializeHandler(context.Background())
	if err != nil {
		panic(err)
	}
	defer cleanup()

	go server.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		fmt.Println("Received signal, shutting down...")
		return
	case err := <-server.Done:
		if err != nil {
			panic(err)
		}
		fmt.Println("Server is done, shutting down...")
		return
	}
}
