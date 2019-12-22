package server

import (
	"context"
	"dubboMesh/agent/dubbo/rpcClient"
	"dubboMesh/agent/dubbo/server/pb"
	"dubboMesh/agent/register/etcdRegister"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

var etcdHost = flag.String("etcd-host", "localhost", "")
var etcdPort = flag.Int("etcd-port", 2379, "")

type Server struct {

}

func (s *Server) Server(ctx context.Context, req *message.AgentRequest) (*message.AgentResponse, error) {
	go func() {
		etcdManager := &etcdRegister.EtcdManager{
			Host: *etcdHost,
			Port: *etcdPort,
		}
		etcdManager.Register(req.Interface, 8000)
	}()

	rClient := new(rpcClient.RpcClient)
	res, err := rClient.Invoke(req.Interface, req.Method, req.ParameterTypesString, req.Parameter)
	if err != nil {
		return nil, err
	}
	resp := &message.AgentResponse{
		RequestID: req.RequestID,
		RespLen: int64(len(res)),
	}
	return resp, nil
}


func Main() {
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		return
	}
	server := new(Server)
	grpcServer := grpc.NewServer()
	message.RegisterAgentServiceServer(grpcServer, server)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")
			grpcServer.GracefulStop()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	log.Println("Listening and serving HTTP on :", port)
	grpcServer.Serve(listen)

}
