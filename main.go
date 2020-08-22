package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gojobprocessor/grpc_service"
	pbm "github.com/gojobprocessor/proto/master"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

var (
	nodeInfo = NodeInfo{}
)

type NodeInfo struct {
	isMaster          *bool
	Port              *int
	MasterHost        *string
	MasterPort        *int
	HeartbeatInterval *int
	ChannelErr        chan<- interface{}
}

func main() {
	nodeInfo.isMaster = flag.Bool("is_master", true, "If true, the node will start as Master, "+
		"else will start as Worker")
	nodeInfo.Port = flag.Int("port", 50001, "This node's port, by default is 50000")
	nodeInfo.MasterHost = flag.String("master_host", "127.0.0.1", "This is the master node's"+
		"host address")
	nodeInfo.MasterPort = flag.Int("master_port", 50000, "This is the master node's port")
	nodeInfo.HeartbeatInterval = flag.Int("heartbeat_int", 5, "If this node is worker, set the "+
		"heartbeat interval, by default 5 seconds")
	flag.Parse()
	if *nodeInfo.isMaster {
		initMaster(nodeInfo)
	} else {
		go initWorker(nodeInfo)
		if err := registerNode(); err != nil {
			log.Fatalf("Could not register against node: %s\n", err)
		}
		sendHeartbeat()
	}
}

func initMaster(nodeInfo NodeInfo) {
	log.Println("Starting node as master...")
	server := &grpc_service.MasterServer{
		Server: grpc_service.Server{Port: strconv.Itoa(*nodeInfo.MasterPort), Network: "tcp"},
	}
	server.ExecGrpcServer()
}

func initWorker(nodeInfo NodeInfo) {
	log.Println("Starting node as worker...")
	server := &grpc_service.WorkerServer{
		Server: grpc_service.Server{Port: strconv.Itoa(*nodeInfo.Port), Network: "tcp"},
	}
	server.ExecGrpcServer()
}

func dial(nodeInfo *NodeInfo) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *nodeInfo.MasterHost, *nodeInfo.MasterPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not send heartbeat against master node: %s", err)
	}
	return conn, nil
}

func newClient(conn *grpc.ClientConn) (pbm.JobProcessorClient, context.Context, context.CancelFunc) {
	client := pbm.NewJobProcessorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return client, ctx, cancel
}

func registerNode() error {
	conn, err := dial(&nodeInfo)
	defer conn.Close()
	client, ctx, cancel := newClient(conn)
	defer cancel()
	wi := pbm.WorkerInfo{
		Host: *nodeInfo.MasterHost,
		Port: int32(*nodeInfo.MasterPort),
	}
	_, err = client.RegisterWorker(ctx, &wi)
	if err != nil {
		log.Printf("Could not send Heartbeat: %s\n", err)
	}
	return nil
}

func sendHeartbeat() {
	for {
		log.Println("Sending heartbeat...")
		conn, err := dial(&nodeInfo)
		client, ctx, cancel := newClient(conn)
		_, err = client.Heartbeat(ctx, &empty.Empty{})
		if err != nil {
			log.Printf("Could not send Heartbeat: %s\n", err)
		}
		cancel()
		conn.Close()
		time.Sleep(time.Duration(*nodeInfo.HeartbeatInterval) * time.Second)
	}

}
