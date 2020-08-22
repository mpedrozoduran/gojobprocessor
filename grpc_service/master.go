package grpc_service

import (
	"context"
	"fmt"
	"github.com/gojobprocessor-master/proto/master"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (s *MasterServer) GetJob(ctx context.Context, in *master.JobId) (*master.JobStatus, error) {
	log.Printf("Calling GetJob: %v", in)
	return &master.JobStatus{}, nil
}

func (s *MasterServer) RunJob(ctx context.Context, in *master.JobRequest) (*master.JobStatus, error) {
	log.Printf("Calling RunJob: %v", in)
	return &master.JobStatus{}, nil
}

func (s *MasterServer) CancelJob(ctx context.Context, in *master.JobId) (*master.JobStatus, error) {
	log.Printf("Calling CancelJob: %v", in)
	return &master.JobStatus{}, nil
}

func (s *MasterServer) RemoveJob(ctx context.Context, in *master.JobId) (*master.JobStatus, error) {
	log.Printf("Calling RemoveJob: %v", in)
	return &master.JobStatus{}, nil
}

func (s *MasterServer) RegisterWorker(ctx context.Context, in *master.WorkerInfo) (*master.RegisterStatus, error) {
	log.Printf("Calling RegisterWorker: %v", in)
	return &master.RegisterStatus{}, nil
}

func (s *MasterServer) Heartbeat(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	log.Printf("Calling Heartbeat: %v", in)
	return &empty.Empty{}, nil
}

func (s *MasterServer) ExecGrpcServer() {
	lis, err := net.Listen(s.Network, fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		log.Fatalf("Could not start master node: %s", err)
	}
	srv := grpc.NewServer()
	master.RegisterJobProcessorServer(srv, s)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
