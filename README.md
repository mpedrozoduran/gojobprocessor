### Re-generate gRPC code
`protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative proto/jobproc.proto`

### Pending changes
- Implement TLS on gRPC client/server