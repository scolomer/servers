package main

import (
	"context"
	"log"
	"net"
	"servers/grpc/gen/service/test"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	keys map[string]string = make(map[string]string)
	mu   sync.Mutex
)

type server struct {
	test.UnimplementedTestServer
}

func (s *server) Test(ctx context.Context, in *test.Request) (*test.Response, error) {
	mu.Lock()
	defer mu.Unlock()

	keys[in.Key] = in.Value
	return &test.Response{}, nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	port := 10101

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		zap.L().Sugar().Fatalf("Unable to bind to port : %v", err)
	}
	s := grpc.NewServer()
	test.RegisterTestServer(s, &server{})
	zap.L().Sugar().Infof("gRPC started on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error : %v", err)
	}
}
