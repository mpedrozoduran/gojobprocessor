package worker

import (
	"context"
	"fmt"
	server "github.com/mpedrozoduran/gojobprocessor/grpc_service"
	"github.com/mpedrozoduran/gojobprocessor/proto/worker"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ServerWorker struct {
	worker.UnimplementedJobProcessorServer
	server.Server
}

func (s *ServerWorker) GetJob(ctx context.Context, in *worker.JobId) (*worker.JobStatus, error) {
	log.Printf("Calling GetJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *ServerWorker) RunJob(ctx context.Context, in *worker.JobRequest) (*worker.JobStatus, error) {
	log.Printf("Calling GetJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *ServerWorker) CancelJob(ctx context.Context, in *worker.JobId) (*worker.JobStatus, error) {
	log.Printf("Calling CancelJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *ServerWorker) RemoveJob(ctx context.Context, in *worker.JobId) (*worker.JobStatus, error) {
	log.Printf("Calling RemoveJob: %v", in)
	return &worker.JobStatus{}, nil
}

func (s *ServerWorker) ExecGrpcServer() {
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
