package grpc_service

import (
	"github.com/gojobprocessor-master/proto/master"
	"github.com/gojobprocessor-master/proto/worker"
)

type Server struct {
	Host    string
	Port    string
	Network string
}

type MasterServer struct {
	master.UnimplementedJobProcessorServer
	Server
}

type WorkerServer struct {
	worker.UnimplementedJobProcessorServer
	Server
}
