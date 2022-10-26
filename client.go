package main

import (
	"bufio"
	"chatServer/protoFiles"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	clientProfile := clientProfile{stream: stream, connection: *conn}
	startScreen()
	clientProfile.chooseUserName()
	go clientProfile.receiveMessage()
	go clientProfile.sendMessage()

	bl := make(chan bool)
	<-bl
}

type clientProfile struct {
	stream     protoFiles.Services_ChatServiceClient
	clientName string
	connection grpc.ClientConn
}

func (cp *clientProfile) chooseUserName() {
	Scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Choose your username: ")
	Scanner.Scan()
	cp.clientName = Scanner.Text()
	cp.joinChatService()

}

func (cp *clientProfile) sendMessage() {

	for {
		Scanner := bufio.NewScanner(os.Stdin)
		Scanner.Scan()
		MessageToBeSent := Scanner.Text()
		if strings.ToLower(MessageToBeSent) == "exit" {

			MessageFromClient := &protoFiles.FromClient{
				Name: cp.clientName,
				Body: "-- has left the chat --",
			}
			err := cp.stream.Send(MessageFromClient)
			if err != nil {
				log.Printf("Failed to send exit-message to server: %s", err)
			}
			time.Sleep(1 * time.Second)
			cp.exitChatService()
		}
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
		fmt.Println()

	}
}

func (cp *clientProfile) exitChatService() {
	cp.connection.Close()
	os.Exit(0)
}

func (cp *clientProfile) joinChatService() {
	MessageFromClient := &protoFiles.FromClient{
		Name: cp.clientName,
		Body: "-- has joined the chat --",
	}
	err := cp.stream.Send(MessageFromClient)
	if err != nil {
		log.Printf("Failed to send join-message: %s", err)
	}
}

func startScreen() {
	fmt.Println("░██╗░░░░░░░██╗███████╗██╗░░░░░░█████╗░░█████╗░███╗░░░███╗███████╗  ████████╗░█████╗░")
	fmt.Println("░██║░░██╗░░██║██╔════╝██║░░░░░██╔══██╗██╔══██╗████╗░████║██╔════╝  ╚══██╔══╝██╔══██╗")
	fmt.Println("░╚██╗████╗██╔╝█████╗░░██║░░░░░██║░░╚═╝██║░░██║██╔████╔██║█████╗░░  ░░░██║░░░██║░░██║")
	fmt.Println("░░████╔═████║░██╔══╝░░██║░░░░░██║░░██╗██║░░██║██║╚██╔╝██║██╔══╝░░  ░░░██║░░░██║░░██║")
	fmt.Println("░░╚██╔╝░╚██╔╝░███████╗███████╗╚█████╔╝╚█████╔╝██║░╚═╝░██║███████╗  ░░░██║░░░╚█████╔╝")
	fmt.Println("░░░╚═╝░░░╚═╝░░╚══════╝╚══════╝░╚════╝░░╚════╝░╚═╝░░░░░╚═╝╚══════╝  ░░░╚═╝░░░░╚════╝░")

	fmt.Println("░█████╗░██╗░░██╗██╗████████╗░█████╗░██╗░░██╗░█████╗░████████╗██╗")
	fmt.Println("██╔══██╗██║░░██║██║╚══██╔══╝██╔══██╗██║░░██║██╔══██╗╚══██╔══╝██║")
	fmt.Println("██║░░╚═╝███████║██║░░░██║░░░██║░░╚═╝███████║███████║░░░██║░░░██║")
	fmt.Println("██║░░██╗██╔══██║██║░░░██║░░░██║░░██╗██╔══██║██╔══██║░░░██║░░░╚═╝")
	fmt.Println("╚█████╔╝██║░░██║██║░░░██║░░░╚█████╔╝██║░░██║██║░░██║░░░██║░░░██╗")
	fmt.Println("░╚════╝░╚═╝░░╚═╝╚═╝░░░╚═╝░░░░╚════╝░╚═╝░░╚═╝╚═╝░░╚═╝░░░╚═╝░░░╚═╝")
}
