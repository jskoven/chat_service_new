package main

import (
	"log"
	"net"
	"os"
	"time"

	"chatServer/protoFiles"

	"google.golang.org/grpc"
)

func main() {
	f, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	log.SetOutput(f)
	listener, err := net.Listen("tcp", ":9000")
	log.Println()
	log.Println()
	log.Printf("### Logs for chat session started at: %s ###", time.Now())
	log.Printf("Starting server...")

	if err != nil {
		log.Printf("Failed to listen on port :9000, Error: %s", err)
	}

	log.Printf("Listening on port :9000")

	grpcserver := grpc.NewServer()

	server := protoFiles.ChatServer{}
	protoFiles.RegisterServicesServer(grpcserver, &server)

	err = grpcserver.Serve(listener)
	if err != nil {
		log.Printf("Failed to serve with listener, error: %s", err)
	}

}
