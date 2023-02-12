### Re-generate gRPC code
  ```bash
  brew install protobuf
  brew install protoc-gen-go-grpc
  ```

  ```bash
  protoc --go_out=. --go-grpc_out=. proto/jobproc-master.proto
  protoc --go_out=. --go-grpc_out=. proto/jobproc-worker.proto
  ```

### Pending changes
- Implement TLS on gRPC client/server

### Ideas
- Types of Jobs: 
  - Large database queries
  - Execute scripts, for example Python