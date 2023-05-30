package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
)

func main() {
	flag.IntVar(&Port, "port", 8081, "set the listening port")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := startTaskBot(ctx)
	if err != nil {
		panic(err)
	}
}
