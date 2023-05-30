package main

import (
	context "context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func AccessStreamInterceptor(logSubs *EventSubs, stats *StatTracker) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		event := getEvent(ss.Context(), info.FullMethod)
		logSubs.Publish(event)
		stats.Track(event.Method, event.Consumer)

		return handler(srv, ss)
	}
}

func AccessUnaryInterceptor(logSubs *EventSubs, stats *StatTracker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		event := getEvent(ctx, info.FullMethod)
		logSubs.Publish(event)
		stats.Track(event.Method, event.Consumer)

		return handler(ctx, req)
	}
}

func AuthStreamInterceptor(ACL map[string][]string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		consumer := getConsumer(ss.Context())
		if !checkAccess(info.FullMethod, ACL[consumer]) {
			return status.Errorf(codes.Unauthenticated, "consumer '%s' does not have access to %s", consumer, info.FullMethod)
		}

		return handler(srv, ss)
	}
}

func AuthUnaryInterceptor(ACL map[string][]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		consumer := getConsumer(ctx)
		if !checkAccess(info.FullMethod, ACL[consumer]) {
			return nil, status.Errorf(codes.Unauthenticated, "consumer '%s' does not have access to %s", consumer, info.FullMethod)
		}

		return handler(ctx, req)
	}
}

func checkAccess(method string, ACL []string) bool {
	for _, m := range ACL {
		if method == m {
			return true
		}
		if strings.HasSuffix(m, "*") && strings.HasPrefix(method, m[:len(m)-1]) {
			return true
		}
	}
	return false
}

func getEvent(ctx context.Context, method string) *Event {
	host := ""
	if p, ok := peer.FromContext(ctx); ok {
		host = p.Addr.String()
	}

	event := &Event{
		Timestamp: time.Now().Unix(),
		Consumer:  getConsumer(ctx),
		Method:    method,
		Host:      host,
	}
	return event
}

func getConsumer(ctx context.Context) string {
	consumer := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		v := md.Get("consumer")
		if len(v) != 0 {
			consumer = v[0]
		}
	}
	return consumer
}
