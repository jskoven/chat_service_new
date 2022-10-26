package main

import (
	"bufio"
	"chatServer/protoFiles"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var clientLamport = 0

func main() {
	f, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	log.SetOutput(f)
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
	stream      protoFiles.Services_ChatServiceClient
	clientName  string
	connection  grpc.ClientConn
	lamportTime int32
}

func (cp *clientProfile) chooseUserName() {
	Scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("To exit the chat service, simply type: exit")
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
		if MessageToBeSent == "myLamport" {
			fmt.Printf("Current lamport: %d", cp.lamportTime)
			continue
		}
		cp.lamportTime++
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
			Name:    cp.clientName,
			Body:    MessageToBeSent,
			Lamport: cp.lamportTime,
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
		cp.lamportTime = int32(math.Max(float64(cp.lamportTime), float64(messageReceived.Lamport)))
		cp.lamportTime = cp.lamportTime + 1
		if err != nil {
			log.Printf("Failed to receive message from server: %s", err)
		}
		fmt.Printf("## Lamport time %d ##   %s : %v", cp.lamportTime, messageReceived.Name, messageReceived.Body)
		fmt.Println()

	}
}

func (cp *clientProfile) exitChatService() {
	log.Printf("User has exited with username: %s", cp.clientName)
	cp.connection.Close()
	os.Exit(0)
}

func (cp *clientProfile) joinChatService() {
	MessageFromClient := &protoFiles.FromClient{
		Name: cp.clientName,
		Body: "-- has joined the chat --",
	}
	log.Printf("User has connected with username: %s", cp.clientName)
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
