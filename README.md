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

## License

[MIT](https://github.com/iamrajiv/opentelemetry-grpc-gateway-boilerplate/blob/main/LICENSE)
