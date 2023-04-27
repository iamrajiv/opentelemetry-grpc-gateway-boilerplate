package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/proto/helloworld/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.GreeterServiceSayHelloRequest) (*pb.GreeterServiceSayHelloResponse, error) {
	return &pb.GreeterServiceSayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(nil)))
	pb.RegisterGreeterServiceServer(s, &server{})

	reflection.Register(s)

	fmt.Println("Starting gRPC server on :8080...")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

func runRESTServer() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(nil))}
	if err := pb.RegisterGreeterServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts); err != nil {
		return err
	}

	fmt.Println("Starting gRPC-Gateway server on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		return err
	}

	return nil
}

func main() {
	go func() {
		if err := runRESTServer(); err != nil {
			panic(err)
		}
	}()

	if err := runGRPCServer(); err != nil {
		panic(err)
	}
}
