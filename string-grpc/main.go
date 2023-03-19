package main

import (
	"context"
	"fmt"
	"github.com/goroo-o56/go-grpcsrv-v2/string-grpc/health"
	examplepb "github.com/goroo-o56/go-grpcsrv-v2/string-grpc/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	examplepb.RegisterStringServiceServer(s, NewServer())
	grpc_health_v1.RegisterHealthServer(s, &health.Server{})
	reflection.Register(s) //grpcurl用にリフレクションする。

	go func() {
		log.Printf("start gRPC server port: %v", port)
		err := s.Serve(listener)
		if err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()

}

func NewServer() *StringServiceServer {
	//これはコンストラクタ。呼び出すことで実体が返される
	return &StringServiceServer{}
}

type StringServiceServer struct {
	examplepb.UnimplementedStringServiceServer
}

func (s *StringServiceServer) ProcessStrings(ctx context.Context, strArr *examplepb.StringArray) (*examplepb.StringResult, error) {
	str := strings.Join(strArr.Values, "")
	return &examplepb.StringResult{Value: str}, nil
}
