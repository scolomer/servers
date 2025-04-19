package main

import (
	"context"
	"math/rand"
	"servers/grpc/gen/service/test"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	port := 10101

	conn, err := grpc.NewClient("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Sugar().Fatalf("did not connect: %v", err)
	}

	c := test.NewTestClient(conn)

	count := 100000

	start := time.Now()
	for i := 0; i < count; i++ {
		_, err := c.Test(context.TODO(), &test.Request{Key: randString(16), Value: "bonjour"})
		if err != nil {
			zap.L().Sugar().Fatalf("Error: %v", err)
		}

		if i%1000 == 0 {
			zap.L().Sugar().Infof("gRPC %v", i)
		}
	}
	end := time.Now()
	tp := float64(count) / end.Sub(start).Seconds()
	zap.L().Sugar().Infof("throughput %v", tp)

}
