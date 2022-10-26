package protoFiles

import (
	"log"
	"math/rand"
	"sync"
)

type messageStruct struct {
	clientName  string
	messageBody string
	messageCode int
	clientCode  int
	lamport     int
}

type messageHandler struct {
	MessageSlice []messageStruct
	lock         sync.Mutex
}

var messageHandlerObject = messageHandler{}

type ChatServer struct {
	UnimplementedServicesServer
}

// Handles both receiving and sending messages
func (s *ChatServer) ChatService(c Services_ChatServiceServer) error {
	clientCode := int(rand.Int31())
	errorChannel := make(chan error)

	//receive message
	go receiveFromStream(c, int(clientCode), errorChannel)

	//send message
	go sendToStream(c, int(clientCode), errorChannel)

	return <-errorChannel
}

func receiveFromStream(c Services_ChatServiceServer, clientCodeReceived int, errorChannel chan error) {
	for {
		message, err := c.Recv()
		if err != nil {
			log.Printf("Failed to receive message from client, error: %s", err)
			errorChannel <- err
		} else {

			messageHandlerObject.lock.Lock()

			messageHandlerObject.MessageSlice = append(messageHandlerObject.MessageSlice, messageStruct{
				clientName:  message.Name,
				messageBody: message.Body,
				messageCode: int(rand.Int31()),
				clientCode:  clientCodeReceived,
				lamport:     int(message.Lamport),
			})
			lastMessage := messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1]
			log.Printf("Broadcasting message from %s. Messagebody: %s. Messagecode: %d. Clientcode: %d. Lamport: %d", lastMessage.clientName, lastMessage.messageBody, lastMessage.messageCode, lastMessage.clientCode, lastMessage.lamport)
			messageHandlerObject.lock.Unlock()

		}

	}
}

// Server sending received messages out to other clients
func sendToStream(c Services_ChatServiceServer, clientCodeSent int, errorChannel chan error) {
	var messageCode int
	for {

		for {
			messageHandlerObject.lock.Lock()

			if len(messageHandlerObject.MessageSlice) == 0 {
				messageHandlerObject.lock.Unlock()
				break
			}

			senderCode := messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1].clientCode
			senderName := messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1].clientName
			senderMessage := messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1].messageBody
			messageHandlerObject.lock.Unlock()

			if senderCode != clientCodeSent && messageCode != messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1].messageCode {
				messageCode = messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1].messageCode
				senderLamport := messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1].lamport

				err := c.Send(&FromServer{
					Name:    senderName,
					Body:    senderMessage,
					Lamport: int32(senderLamport),
				})
				if err != nil {
					errorChannel <- err
				}

			}

		}
	}
}
