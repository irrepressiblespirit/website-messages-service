package grpcapi

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/irrepressiblespirit/website-messages-service/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UnaryAuthInterceptor struct {
	Auth    service.IAuthToken
	PrintMD bool
}

func (u *UnaryAuthInterceptor) Interceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("UnaryAuthInterceptor: Error parsing metadata %s\n", info.FullMethod)
		return handler(ctx, req)
	}

	if u.PrintMD {
		log.Printf("UnaryAuthInterceptor %s [%s] Metadata:\n", info.FullMethod, strings.Join(md[":authority"], ","))
		for k, v := range md {
			fmt.Printf("\t%s\t%s\n", k, v[0])
		}
	}

	token, ok := md["token"]
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Token not found")
	}
	if err := u.Auth.Check(token[0]); err != nil {
		log.Printf("UnaryAuthInterceptor Error: %s\n", err)
		return nil, status.Error(codes.PermissionDenied, "Token invalid")
	}

	return handler(ctx, req)
}
