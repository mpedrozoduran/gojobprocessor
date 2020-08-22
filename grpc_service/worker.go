package grpc_service

import (
	"context"
	"fmt"
	"github.com/gojobprocessor-master/proto/worker"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (s *WorkerServer) GetJob(ctx context.Context, in *worker.JobId) (*worker.JobStatus, error) {
	log.Printf("Calling GetJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *WorkerServer) RunJob(ctx context.Context, in *worker.JobRequest) (*worker.JobStatus, error) {
	log.Printf("Calling GetJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *WorkerServer) CancelJob(ctx context.Context, in *worker.JobId) (*worker.JobStatus, error) {
	log.Printf("Calling CancelJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *WorkerServer) RemoveJob(ctx context.Context, in *worker.JobId) (*worker.JobStatus, error) {
	log.Printf("Calling RemoveJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *WorkerServer) ExecGrpcServer() {
	lis, err := net.Listen(s.Network, fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		log.Fatalf("Could not start worker node: %s", err)
	}
	srv := grpc.NewServer()
	worker.RegisterJobProcessorServer(srv, s)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
