package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mpedrozoduran/gojobprocessor/api"
	"github.com/mpedrozoduran/gojobprocessor/grpc_service"
	"github.com/mpedrozoduran/gojobprocessor/master"
	pbm "github.com/mpedrozoduran/gojobprocessor/proto/master"
	"github.com/mpedrozoduran/gojobprocessor/worker"
	"google.golang.org/grpc"
	"log"
	"net/http"
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
	HttpApiPort       *int
	HeartbeatInterval *int
	ChannelErr        chan<- interface{}
}

func main() {
	nodeInfo.isMaster = flag.Bool("is_master", true, "If true, the node will start as ServerMaster, "+
		"else will start as Worker")
	nodeInfo.Port = flag.Int("port", 50001, "This node's port, by default is 50000")
	nodeInfo.MasterHost = flag.String("master_host", "127.0.0.1", "This is the master node's"+
		"host address")
	nodeInfo.MasterPort = flag.Int("master_port", 50000, "This is the master node's port")
	nodeInfo.HttpApiPort = flag.Int("rest_api_port", 8181, "This is the HTTP REST API port")
	nodeInfo.HeartbeatInterval = flag.Int("heartbeat_int", 5, "If this node is worker, set the "+
		"heartbeat interval, by default 5 seconds")
	flag.Parse()
	if *nodeInfo.isMaster {
		go initMaster(nodeInfo)
		initHttpServer(nodeInfo)
	} else {
		go initWorker(nodeInfo)
		worker, err := registerNode()
		if err != nil {
			log.Fatalf("Could not register against node: %s\n", err)
		}
		done := make(chan int, 1)
		go sendHeartbeat(worker, done)
		<-done
	}
}

func initMaster(nodeInfo NodeInfo) {
	log.Println("starting node as master...")
	server := &master.ServerMaster{
		Server: grpc_service.Server{Port: strconv.Itoa(*nodeInfo.MasterPort), Network: "tcp"},
	}
	server.ExecGrpcServer()
}

func initHttpServer(info NodeInfo) {
	rtr := chi.NewRouter()
	router := api.Router{
		MuxRouter: rtr,
	}
	http.ListenAndServe(fmt.Sprintf(":%v", info.HttpApiPort), router.MuxRouter)
}

func initWorker(nodeInfo NodeInfo) {
	log.Println("Starting node as worker...")
	server := &worker.ServerWorker{
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

func registerNode() (pbm.WorkerInfo, error) {
	conn, err := dial(&nodeInfo)
	defer conn.Close()
	client, ctx, cancel := newClient(conn)
	defer cancel()
	wi := pbm.WorkerInfo{
		Id:    base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", *nodeInfo.MasterHost, int32(*nodeInfo.MasterPort)))),
		Host:  *nodeInfo.MasterHost,
		Port:  int32(*nodeInfo.MasterPort),
		State: pbm.WorkerState_STARTING,
	}
	_, err = client.RegisterWorker(ctx, &wi)
	if err != nil {
		log.Printf("could not register node: %s\n", err)
	}
	return wi, nil
}

func sendHeartbeat(worker pbm.WorkerInfo, done chan<- int) {
	heartbeat := func() error {
		conn, err := dial(&nodeInfo)
		defer conn.Close()
		client, ctx, cancel := newClient(conn)
		defer cancel()
		newWorker := pbm.WorkerInfo{
			Id:    worker.Id,
			Host:  worker.Host,
			Port:  worker.Port,
			Name:  worker.Name,
			State: pbm.WorkerState_UP,
		}
		_, err = client.Heartbeat(ctx, &newWorker)
		if err != nil {
			errNotSent := fmt.Errorf("could not send heartbeat: %s\n", err)
			return errNotSent
		}

		return nil
	}
	maxFailedHeartbeats := 3
	currFailedHeartbeats := 0
	for {
		log.Println("sending heartbeat...")
		time.Sleep(time.Duration(*nodeInfo.HeartbeatInterval) * time.Second)
		err := heartbeat()
		if err == nil {
			currFailedHeartbeats = 0
			continue
		}
		log.Println(err)
		currFailedHeartbeats++
		if maxFailedHeartbeats == currFailedHeartbeats {
			log.Printf("reached maxFailedHeartbeats: %v. Shutting down worker.\n", maxFailedHeartbeats)
			goto end
		}
	}
end:
	done <- 1
}
