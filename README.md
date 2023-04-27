<div align="center">
<img src="assets/opentelemetry-grpc-gateway-boilerplate.svg" height="auto" width="400" />
<br />
<h1>OpenTelemetry gRPC-Gateway Boilerplate</h1>
<p>
Boilerplate for gRPC-Gateway with OpenTelemetry instrumentation

</p>
<a href="https://github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/network/members"><img src="https://img.shields.io/github/forks/iamrajiv/opentelemetry-grpc-gateway-boilerplate?color=0969da&style=for-the-badge" height="auto" width="auto" /></a>
<a href="https://github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/stargazers"><img src="https://img.shields.io/github/stars/iamrajiv/opentelemetry-grpc-gateway-boilerplate?color=0969da&style=for-the-badge" height="auto" width="auto" /></a>
<a href="https://github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/blob/main/LICENSE"><img src="https://img.shields.io/github/license/iamrajiv/opentelemetry-grpc-gateway-boilerplate?color=0969da&style=for-the-badge" height="auto" width="auto" /></a>
</div>

## About

This is an example gRPC and gRPC-Gateway service that implements a "Hello World" RPC and REST API. The service is written in Go and uses the gRPC framework and the gRPC-Gateway library to provide a gRPC service and a RESTful API that maps to the same underlying gRPC service.

#### Folder structure:

```shell
.
├── LICENSE
├── Makefile
├── README.md
├── assets
│   └── opentelemetry-grpc-gateway-boilerplate.svg
├── buf.gen.yaml
├── go.mod
├── go.sum
├── main.go
└── proto
    ├── buf.lock
    ├── buf.yaml
    └── helloworld
        └── v1
            ├── helloworld
            │   └── v1
            │       ├── helloworld.pb.go
            │       ├── helloworld.pb.gw.go
            │       └── helloworld_grpc.pb.go
            └── helloworld.proto
```

## Usage

#### Generating Protobuf and gRPC-Gateway Code

To generate the protobuf and gRPC-Gateway code, you can use the buf tool. The protobuf files are stored in the proto directory, and the generated code will be placed in the internal directory.

To generate the code, run the following command:

```shell
make generate
```

This will generate the necessary Go code for the protobuf files and the gRPC-Gateway files.

#### gRPC Requests and Responses

Once you've started both the gRPC and gRPC-Gateway servers using `go run main.go`, you can send gRPC requests and receive responses using a gRPC client.

Here's an example of how to send a gRPC request and receive a response using the `grpcurl` command-line tool:

1. Install `grpcurl` by following the instructions in the [official documentation](https://github.com/fullstorydev/grpcurl#installation).
2. Open a new terminal window or tab and run the following command to send a gRPC request to the `SayHello` RPC:

```shell
grpcurl -plaintext -d '{"name": "John"}' localhost:8080 helloworld.v1.GreeterService/SayHello
```

This sends a gRPC request to the `SayHello` RPC with the name "John" as a request parameter.

3. You should receive a response that looks like this:

```shell
{
  "message": "Hello, John!"
}
```

This is the response message returned by the `SayHello` RPC.

#### REST Requests and Responses

Once you've started both the gRPC and gRPC-Gateway servers using `go run main.go`, you can send REST requests and receive responses using a tool like `curl`.

Here's an example of how to send a REST request and receive a response using `curl`:

1. Open a new terminal window or tab and run the following command to send a REST request to the `/v1/helloworld` endpoint:

```shell
curl -X POST http://localhost:8081/v1/helloworld -H "Content-Type: application/json" -d '{"name": "John"}'
```

This sends a `POST` request to the `/v1/helloworld` endpoint with a JSON payload containing the name "John".

2. You should receive a response that looks like this:

```shell
{
  "message": "Hello, John!"
}
```

#### Adding OpenTelemetry Instrumentation to gRPC-Gateway

To add simple OpenTelemetry instrumentation to your gRPC-Gateway, follow these steps:

1. Import necessary packages:

```go
import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)
```

2. Add a function to initialize OpenTelemetry:

```go
func initTracing() func() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

	// Create a simple span processor that writes to the exporter.
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))
	otel.SetTracerProvider(tp)

	// Set the global propagator to use W3C Trace Context.
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Return a function to stop the tracer provider.
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}
}
```

3. Modify the gRPC-Gateway function to pass the OpenTelemetry context to the gRPC client connection:

```go
func runRESTServer() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	// Use the OpenTelemetry gRPC client interceptor for tracing.
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		return err
	}

	if err := pb.RegisterGreeterServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Println("Starting gRPC-Gateway server on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		return err
	}

	return nil
}
```

3. In the main function, initialize OpenTelemetry and defer the shutdown function:

```go

func main() {
	// Initialize tracing and handle the tracer provider shutdown.
	stopTracing := initTracing()
	defer stopTracing()

	// Start the REST server in a goroutine.
	go func() {
		if err := runRESTServer(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	// Start the gRPC server.
	if err := runGRPCServer(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
```

To test the OpenTelemetry instrumentation we will send gRPC-Gateway requests using `curl` and view the traces. So when I send a request to the `/v1/helloworld` endpoint, I should see a trace in the console that looks like this:

```shell
➜  opentelemetry-grpc-gateway-boilerplate git:(main) ✗ go run main.go
Starting gRPC server on :8080...
Starting gRPC-Gateway server on :8081...
{
	"Name": "helloworld.v1.GreeterService/SayHello",
	"SpanContext": {
		"TraceID": "3a2c50516735224c323ffb848399819d",
		"SpanID": "60263bb1d968fb06",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "00000000000000000000000000000000",
		"SpanID": "0000000000000000",
		"TraceFlags": "00",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 3,
	"StartTime": "2023-04-27T22:28:47.164181+05:30",
	"EndTime": "2023-04-27T22:28:47.164874369+05:30",
	"Attributes": [
		{
			"Key": "rpc.system",
			"Value": {
				"Type": "STRING",
				"Value": "grpc"
			}
		},
		{
			"Key": "rpc.service",
			"Value": {
				"Type": "STRING",
				"Value": "helloworld.v1.GreeterService"
			}
		},
		{
			"Key": "rpc.method",
			"Value": {
				"Type": "STRING",
				"Value": "SayHello"
			}
		},
		{
			"Key": "net.peer.name",
			"Value": {
				"Type": "STRING",
				"Value": "localhost"
			}
		},
		{
			"Key": "net.peer.port",
			"Value": {
				"Type": "INT64",
				"Value": 8080
			}
		},
		{
			"Key": "rpc.grpc.status_code",
			"Value": {
				"Type": "INT64",
				"Value": 0
			}
		}
	],
	"Events": [
		{
			"Name": "message",
			"Attributes": [
				{
					"Key": "message.type",
					"Value": {
						"Type": "STRING",
						"Value": "SENT"
					}
				},
				{
					"Key": "message.id",
					"Value": {
						"Type": "INT64",
						"Value": 1
					}
				}
			],
			"DroppedAttributeCount": 0,
			"Time": "2023-04-27T22:28:47.16422+05:30"
		},
		{
			"Name": "message",
			"Attributes": [
				{
					"Key": "message.type",
					"Value": {
						"Type": "STRING",
						"Value": "RECEIVED"
					}
				},
				{
					"Key": "message.id",
					"Value": {
						"Type": "INT64",
						"Value": 1
					}
				}
			],
			"DroppedAttributeCount": 0,
			"Time": "2023-04-27T22:28:47.164873+05:30"
		}
	],
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:main"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.14.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
		"Version": "semver:0.40.0",
		"SchemaURL": ""
	}
}
```

## License

[MIT](https://github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/blob/main/LICENSE)
