package grpc

import (
	"context"
	"google.golang.org/grpc"
	"gprc/pkg/api/v1"
	"log"
	"net"
	"os"
	"os/signal"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, todoServer v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, todoServer)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
