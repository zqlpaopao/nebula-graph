package server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	v11 "test/grpc/api/proto"
)

func RunServer(ctx context.Context, v1API v11.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if nil != err {
		return err
	}
	server := grpc.NewServer()
	v11.RegisterToDoServiceServer(server, v1API)
	c := make(chan os.Signal, 1)
	go func() {
		for range c {
			log.Println("shutting down GRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	log.Println("start gRPC server...")
	return server.Serve(listen)

}
