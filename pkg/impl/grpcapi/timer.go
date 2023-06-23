package grpcapi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func TimeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ip := ":error"
	start := time.Now()
	methodStatus := "OK"
	methodColor := green
	errText := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ip = strings.Join(md[":authority"], ",")
	} else {
		log.Println("Error TimeInterceptor, can't load metadata")
	}
	data, err := handler(ctx, req)
	if err != nil {
		methodColor = red
		methodStatus = "ERR"
		errText = err.Error() + "\n"
	}
	elapsed := time.Since(start)
	fmt.Printf("[gRPC]%v |%s %-3s %s| %13v | %15s |%s %-7s %s %#v\n%s",
		time.Now().Format("2006/01/02 - 15:04:05"),
		methodColor, methodStatus, reset,
		elapsed,
		ip,
		blue, "", reset,
		info.FullMethod,
		errText,
	)
	return data, err
}
