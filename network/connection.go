package network

import (
	"fmt"
	"math"
	"net"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"tcp.com/config"
	pb "tcp.com/proto"
)

type Connection struct {
	Address       string `json:"address"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	TCPConnection net.Conn
}

func (conn *Connection) SendMessage(message *pb.Message) error {
	body := append([]byte{}, message.Body...)
	maxBodySize := 50000
	totalPackage := math.Ceil(float64(len(body)) / float64(maxBodySize))
	if totalPackage == 0 {
		totalPackage = 1
	}
	message.Header.TotalPackage = int32(totalPackage)
	message.Header.Id = uuid.New().String()
	for i := 0; int32(i) < message.Header.TotalPackage; i++ {
		var sendBody []byte
		if len(body) < maxBodySize {
			sendBody = body
		} else {
			sendBody = body[:maxBodySize]
			body = body[maxBodySize:]
		}
		sendMessage := &pb.Message{
			Header: message.Header,
			Body:   sendBody,
		}

		// b, err := json.Marshal(sendMessage)
		b, err := proto.Marshal(sendMessage)
		if err != nil {
			fmt.Printf("Error when marshal %v", err)
			return err
		}
		conn.TCPConnection.Write(b)
	}
	return nil
}

func (conn *Connection) SendInitConnection() {
	protoRs, _ := proto.Marshal(&pb.InitConnection{
		Address: config.AppConfig.Address,
	})
	message := &pb.Message{
		Header: &pb.Header{
			Type:    "request",
			From:    config.AppConfig.Address,
			Command: "InitConnection",
		},
		Body: protoRs,
	}

	err := conn.SendMessage(message)
	if err != nil {
		fmt.Printf("Error when send started %v", err)
	}
}
