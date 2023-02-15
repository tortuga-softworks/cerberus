# cerberus
Authentication server

## Configuration

Environment variables: 

| Variable                  | Default value | Description                                |
|---------------------------|---------------|--------------------------------------------|
| CERBERUS_PORT             | 9000          | The port the application listens on        |
| CERBERUS_SESSION_DURATION | 43200         | The duration of a user session, in seconds |
| CERBERUS_CACHE_HOST       |               | The host of the cache (Redis)              |
| CERBERUS_CACHE_PORT       |               | The port of the cache (Redis)              |

## API (TODO)
### Examples:

    PS C:\projects\cerberus> grpcurl -plaintext -d '{\"username\": \"marem\"}' localhost:9000 cerberus.Authentication/Login
    {
    "sessionId": "jk0obS-CSywOGgGR74NMrlJP5N-77nH5t10MBgAmHHs="
    }
    PS C:\projects\cerberus> grpcurl -plaintext -d '{\"username\": \"\"}' localhost:9000 cerberus.Authentication/Login     
    ERROR:
    Code: InvalidArgument
    Message: the username cannot be empty

## Build

The application is meant to be built using the provided Dockerfile. However, you can also do it manually.

Requirements:
- Go (v1.20)

Command:
    
    go build -o cerberus ./src/cmd

This is assuming the gRPC files have been generated already. If they are not, please reading the following section.

## Development

### Requirements: 
- Go (v1.20)
- Protocol buffer compiler (v3)
- Go plugins:
    - protoc-gen-go (v1.28)
    - protoc-gen-go-grpc (v1.2)

Note: Make sure protoc (Protocol buffer compiler) can find the plugins in the Path environment variable

See https://grpc.io/docs/languages/go/quickstart/

### Generating gRPC source files
To generate the gRPC server and client source files:
    
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative src/cerberus.proto
