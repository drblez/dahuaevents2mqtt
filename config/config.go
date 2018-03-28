package config

import (
	"os"
	"encoding/json"
	"dahuaevents2mqtt/camera"
	"dahuaevents2mqtt/mqtt"
)

const configFile = "dahuaevents2mqtt.json"

type Configuration struct {
	Defaults camera.Config
	Cameras []camera.Config
	MQTT mqtt.Config
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
	for i := range conf.Cameras {
		if conf.Cameras[i].Map == nil {
			conf.Cameras[i].Map = conf.Defaults.Map
		}
		if conf.Cameras[i].Events == nil {
			conf.Cameras[i].Events = conf.Defaults.Events
		}
		if conf.Cameras[i].Password == "" {
			conf.Cameras[i].Password = conf.Defaults.Password
		}
		if conf.Cameras[i].Username == "" {
			conf.Cameras[i].Username = conf.Defaults.Username
		}
		if conf.Cameras[i].Port == "" {
			conf.Cameras[i].Port = conf.Defaults.Port
		}
	}
	return conf
}