package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"tcp.com/dataType"
)

type IConnection interface{}

type Config struct {
	Address    string               `json:"address"`
	Ip         string               `json:"ip"`
	Port       int                  `json:"port"`
	Validators []dataType.Validator `json:"validators"`
}

func loadConfig() Config {
	var config Config
	raw, err := ioutil.ReadFile("config/conf.json")
	if err != nil {
		log.Fatalf("Error occured while reading config")
	}
	json.Unmarshal(raw, &config)
	log.Printf("Config loaded: %v\n", config)
	return config
}

var AppConfig = loadConfig()
