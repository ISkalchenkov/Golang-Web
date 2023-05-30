package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
// если хочется, то для красоты можно разнести логику по разным файликам

func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	ACL := map[string][]string{}
	err := json.Unmarshal([]byte(ACLData), &ACL)
	if err != nil {
		return fmt.Errorf("unmarshal ACLData failed: %w", err)
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("listening on %s failed: %w", listenAddr, err)
	}

	subs := NewEventSubs()
	stats := NewStatTracker()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			AccessUnaryInterceptor(subs, stats),
			AuthUnaryInterceptor(ACL),
		),
		grpc.ChainStreamInterceptor(
			AccessStreamInterceptor(subs, stats),
			AuthStreamInterceptor(ACL),
		),
	)

	RegisterAdminServer(server, NewAdminModule(subs, stats))
	RegisterBizServer(server, NewBizModule())

	errCh := make(chan error)
	go func() {
		if err = server.Serve(lis); err != nil {
			errCh <- fmt.Errorf("grpc server could not serve: %w", err)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			server.GracefulStop()
		case err := <-errCh:
			fmt.Println(err)
		}
	}()

	return nil
}
