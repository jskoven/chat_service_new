package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"github.com/jskoven/chat_service_new/protoFiles"
)

func main() {

	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		log.Printf("Failed to listen on port :9000, Error: %s", err)
	}

	log.Printf("Listening on port :9000")

	grpcserver := grpc.NewServer()

	err = grpcserver.Serve(listener)
	if err != nil {
		log.Printf("Failed to serve with listener, error: %s", err)
	}

}
