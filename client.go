package main

import (
	"bufio"
	"chatServer/protoFiles"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial server on port 9000: %s", err)
	}
	defer conn.Close()

	client := protoFiles.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Printf("Failed to create stream to chatService: %s", err)
	}

	clientProfile := clientProfile{stream: stream}
	clientProfile.chooseUserName()
	go clientProfile.receiveMessage()
	go clientProfile.sendMessage()

	bl := make(chan bool)
	<-bl
}

type clientProfile struct {
	stream     protoFiles.Services_ChatServiceClient
	clientName string
}

func (cp *clientProfile) chooseUserName() {
	Scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Choose your username: ")
	Scanner.Scan()
	cp.clientName = Scanner.Text()
}

func (cp *clientProfile) sendMessage() {

	for {
		Scanner := bufio.NewScanner(os.Stdin)
		Scanner.Scan()
		MessageToBeSent := Scanner.Text()
		MessageFromClient := &protoFiles.FromClient{
			Name: cp.clientName,
			Body: MessageToBeSent,
		}
		err := cp.stream.Send(MessageFromClient)
		if err != nil {
			log.Printf("Failed to send message to server: %s", err)
		}
	}
}

func (cp *clientProfile) receiveMessage() {
	for {
		messageReceived, err := cp.stream.Recv()
		if err != nil {
			log.Printf("Failed to receive message from server: %s", err)
		}
		fmt.Printf("%s : %v", messageReceived.Name, messageReceived.Body)

	}
}
