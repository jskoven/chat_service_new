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
}

type messageHandler struct {
	MessageSlice []messageStruct
	lock         sync.Mutex
}

var messageHandlerObject = messageHandler{}

type chatServer struct {
}

// Handles both receiving and sending messages
func (s *chatServer) chatService(c Services_ChatServiceServer) error {
	clientCode := rand.Int31n(10)
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
			})
			messageHandlerObject.lock.Unlock()
			log.Printf("%v", messageHandlerObject.MessageSlice[len(messageHandlerObject.MessageSlice)-1])
		}

	}
}

// Server sending received messages out to other clients
func sendToStream(c Services_ChatServiceServer, clientCodeSent int, errorChannel chan error) {
	for {

		for {
			messageHandlerObject.lock.Lock()

			senderCode := messageHandlerObject.MessageSlice[0].clientCode
			senderName := messageHandlerObject.MessageSlice[0].clientName
			senderMessage := messageHandlerObject.MessageSlice[0].messageBody

			messageHandlerObject.lock.Unlock()

			if senderCode != clientCodeSent {
				err := c.Send(&FromServer{
					Name: senderName,
					Body: senderMessage,
				})
				if err != nil {
					errorChannel <- err
				}

				messageHandlerObject.lock.Lock()

				if len(messageHandlerObject.MessageSlice) > 1 {
					// The ":1" specifies that the slice should be the same
					//slice, but with only the values of lower bound index 1
					//hence, practically deleting the message at index 0.
					messageHandlerObject.MessageSlice = messageHandlerObject.MessageSlice[1:]
				} else {
					messageHandlerObject.MessageSlice = []messageStruct{}
				}
				messageHandlerObject.lock.Unlock()
			}

		}
	}
}
