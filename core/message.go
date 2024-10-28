package core

import (
	"sync"
	"math/rand"
	"chat_service/api"
	"log"
	"io"
	"time"
)

type client struct {
	Name string
	Id int64
}

type message struct {
	Body string
	Id int64
	Client client
}

type messageHandle struct {
	MQue []message
	mu sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
}

func ReceiveFromStream(stream api.ChatService_HandleCommunicationServer, clientUniqueCode int64, errch chan error) {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
				errch <- err
		} else
		if err != nil {
				log.Printf("Error receiving message: %v", err)
				errch <- err
		} else {
			
			messageHandleObject.mu.Lock()
			messageHandleObject.MQue = append(messageHandleObject.MQue, message{
				Body:       msg.Body,
				Id: int64(rand.Intn(1e8)),
				Client: client{
					Name: msg.SenderName,
					Id: int64(clientUniqueCode),
				},
			})
			log.Printf("Received message from %s: %s", msg.SenderName, msg.Body)
			messageHandleObject.mu.Unlock()
		}
	}
}

func SendToStream(stream api.ChatService_HandleCommunicationServer, clientId int64, errch chan error) {
	for {
		for {
			time.Sleep(500 * time.Millisecond)

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderId := messageHandleObject.MQue[0].Client.Id
			senderName4Client := messageHandleObject.MQue[0].Client.Name
			message4Client := messageHandleObject.MQue[0].Body

			messageHandleObject.mu.Unlock()

			if senderId != clientId {

				err := stream.Send(&api.ChatMessage{ReceiverName: senderName4Client, Body: message4Client})

				if err != nil {
					errch <- err
				}

				messageHandleObject.mu.Lock()

				if len(messageHandleObject.MQue) > 1 {
					messageHandleObject.MQue = messageHandleObject.MQue[1:]
				} else {
					messageHandleObject.MQue = []message{}
				}

				messageHandleObject.mu.Unlock()

			}

		}
	}
}