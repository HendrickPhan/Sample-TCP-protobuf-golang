package main

import (
	"fmt"

	"tcp.com/config"
	"tcp.com/network"
)

func initValidatorConnections() []network.Connection {
	validators := config.AppConfig.Validators
	var validatorConnections []network.Connection
	for _, validator := range validators {
		validatorConnections = append(validatorConnections, network.Connection{
			IP:      validator.Ip,
			Port:    validator.Port,
			Address: validator.Address,
		})
	}
	return validatorConnections
}

func main() {
	finish := make(chan bool)
	initedConnectionsChan := make(chan network.Connection)
	removeConnectionChan := make(chan network.Connection)

	messageHandler := network.MessageHandler{
		InitedConnectionsChan: initedConnectionsChan,
		RemoveConnectionChan:  removeConnectionChan,
	}
	server := network.Server{
		MessageHandler:        messageHandler,
		InitedConnections:     make(map[string]network.Connection),
		IP:                    config.AppConfig.Ip,
		Port:                  config.AppConfig.Port,
		Address:               config.AppConfig.Address,
		InitedConnectionsChan: initedConnectionsChan,
		RemoveConnectionChan:  removeConnectionChan,
	}

	go server.Run(initValidatorConnections())

	<-finish
	fmt.Printf("Main\n")
}
