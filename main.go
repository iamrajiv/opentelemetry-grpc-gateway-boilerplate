package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	// Import the required packages
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/proto/helloworld/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Define the server struct, which implements the pb.UnimplementedGreeterServiceServer interface
type server struct {
	pb.UnimplementedGreeterServiceServer
}

// Initialize OpenTelemetry tracing and return a function to stop the tracer provider
func initTracing() func() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

	// Create a simple span processor that writes to the exporter
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))
	otel.SetTracerProvider(tp)

	// Set the global propagator to use W3C Trace Context
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Return a function to stop the tracer provider
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}
}

// Implement the SayHello method of the pb.GreeterServiceServer interface
func (s *server) SayHello(ctx context.Context, req *pb.GreeterServiceSayHelloRequest) (*pb.GreeterServiceSayHelloResponse, error) {
	return &pb.GreeterServiceSayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

// Set up the gRPC server on port 8080 and serve requests indefinitely
func runGRPCServer() error {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServiceServer(s, &server{})

	// Enable reflection to allow clients to query the server's services
	reflection.Register(s)

	fmt.Println("Starting gRPC server on :8080...")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

// Set up the REST server on port 8081 and handle requests by proxying them to the gRPC server
func runRESTServer() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	// Use the OpenTelemetry gRPC client interceptor for tracing
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		return err
	}

	// Register the gRPC server's handler with the HTTP mux
	if err := pb.RegisterGreeterServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Println("Starting gRPC-Gateway server on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		return err
	}
	return nil
}

func main() {
	// Initialize tracing and handle the tracer provider shutdown
	stopTracing := initTracing()
	defer stopTracing()
	// Start the REST server in a goroutine
	go func() {
		if err := runRESTServer(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	// Start the gRPC server
	if err := runGRPCServer(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
