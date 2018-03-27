package config

import (
	"os"
	"encoding/json"
	"dahuaevents2mqtt/camera"
)

const configFile = "dahuaevents2mqtt.json"

type MQTTConfig struct {
	Host string
	Port string
	Timeout int
}

type Configuration struct {
	Cameras []camera.Config
	MQTT MQTTConfig
}

func Init() Configuration {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
	return conf
}