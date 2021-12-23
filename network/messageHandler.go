package network

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	pb "tcp.com/proto"
)

type MessageHandler struct {
	InitedConnectionsChan chan Connection
	RemoveConnectionChan  chan Connection
}

func (handler *MessageHandler) OnConnect(conn Connection) {
	log.Info(fmt.Sprintf("OnConnect with server %s", conn.TCPConnection.RemoteAddr()))
	conn.SendInitConnection()
}

func (handler *MessageHandler) OnDisconnect(conn Connection) {
	log.Info(fmt.Printf("Disconnected with server  %s, wallet address: %v", conn.TCPConnection.RemoteAddr(), conn.Address))
	handler.RemoveConnectionChan <- conn
	// TODO remove from connection list

}

func (handler *MessageHandler) HandleConnection(conn Connection) {
	pendingMessages := make(map[string]*pb.Message) // map between id and message

	for {
		data := make([]byte, 65535)
		length, err := conn.TCPConnection.Read(data)

		if err != nil {
			switch err {
			case io.EOF:
				handler.OnDisconnect(conn)
				return
			default:
				log.Error("server error: %v\n", err)
				return
			}
		}

		// rawMessage := strings.TrimSpace(string(netData))
		// if rawMessage == "STOP" {
		// 	break
		// }
		message := &pb.Message{}
		proto.Unmarshal(data[:length], message)
		log.Info(fmt.Sprintf("message %v", message))

		if pendingMessage, ok := pendingMessages[message.Header.Id]; ok {
			//do something here
			pendingMessage.Body = append(pendingMessage.Body, message.Body...)
			pendingMessage.Header.TotalReceived++
			pendingMessages[message.Header.Id] = pendingMessage
			if pendingMessage.Header.TotalPackage == pendingMessage.Header.TotalReceived {
				// create process messsage and remove from pending
				go handler.ProcessMessage(conn, pendingMessage)
				delete(pendingMessages, message.Header.Id)
			}
		} else {
			if message.Header.TotalPackage == 1 {
				go handler.ProcessMessage(conn, message)
			} else {
				message.Header.TotalReceived = 1
				pendingMessages[message.Header.Id] = message
			}
		}
	}
}

func (handler *MessageHandler) ProcessMessage(conn Connection, message *pb.Message) {
	switch message.Header.Command {
	case "InitConnection":
		handler.handleInitConnectionMessage(conn, message)

		// case "ValidatorStarted":
		// 	handler.handlerValidatorStarted(message)
		// case "LeaderTick":
		// 	handler.handlerLeaderTick(message)
		// case "VotedBlock":
		// 	handler.handlerVotedBlock(message)
		// case "VoteLeaderBlock":
		// 	handler.handlerVoteLeaderBlock(message)
		// case "SendCheckedBlock":
		// 	handler.handlerSendCheckedBlock(message)

	}
}

func (handler *MessageHandler) handleInitConnectionMessage(conn Connection, message *pb.Message) {
	log.Info("Receive InitConnection from", conn.TCPConnection.RemoteAddr())
	initConnectionMessage := &pb.InitConnection{}
	proto.Unmarshal([]byte(message.Body), initConnectionMessage)
	conn.Address = initConnectionMessage.Address
	handler.InitedConnectionsChan <- conn
}
