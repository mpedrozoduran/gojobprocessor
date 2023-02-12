package master

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	server "github.com/mpedrozoduran/gojobprocessor/grpc_service"
	"github.com/mpedrozoduran/gojobprocessor/proto/master"
	"github.com/mpedrozoduran/gojobprocessor/worker"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ServerMaster struct {
	master.UnimplementedJobProcessorServer
	server.Server
	worker.RouteTable
}

func (s *ServerMaster) GetJob(ctx context.Context, in *master.JobId) (*master.JobStatus, error) {
	log.Printf("calling GetJob: %v\n", in)
	return &master.JobStatus{}, nil
}

func (s *ServerMaster) RunJob(ctx context.Context, in *master.JobRequest) (*master.JobStatus, error) {
	log.Printf("calling RunJob: %v\n¬", in)
	return &master.JobStatus{}, nil
}

func (s *ServerMaster) CancelJob(ctx context.Context, in *master.JobId) (*master.JobStatus, error) {
	log.Printf("calling CancelJob: %v\n", in)
	return &master.JobStatus{}, nil
}

func (s *ServerMaster) RemoveJob(ctx context.Context, in *master.JobId) (*master.JobStatus, error) {
	log.Printf("calling RemoveJob: %v\n", in)
	return &master.JobStatus{}, nil
}

func (s *ServerMaster) RegisterWorker(ctx context.Context, in *master.WorkerInfo) (*master.RegisterStatus, error) {
	log.Printf("calling RegisterWorker: %v\n", in)
	s.RouteTable.Push(*in)
	return &master.RegisterStatus{
		Status:  0,
		Message: "OK",
	}, nil
}

func (s *ServerMaster) Heartbeat(ctx context.Context, in *master.WorkerInfo) (*empty.Empty, error) {
	log.Printf("calling Heartbeat: %v\n", in)
	if in.GetState() == master.WorkerState_DOWN {
		err := fmt.Errorf("worker %s down, removing from route table\n", in.Id)
		s.RouteTable.Delete(in.Id)
		return &empty.Empty{}, err
	}
	_, err := s.RouteTable.Get(in.GetId())
	if err != nil {
		log.Println(err)
		return &empty.Empty{}, err
	}
	s.Update(*in)
	return &empty.Empty{}, nil
}

func (s *ServerMaster) ExecGrpcServer() {
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
