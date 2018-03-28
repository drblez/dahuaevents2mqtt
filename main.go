package main

import (
	"dahuaevents2mqtt/camera"
	"dahuaevents2mqtt/config"
	"fmt"
	"dahuaevents2mqtt/event"
	MQTT "dahuaevents2mqtt/mqtt"
)

func main() {

	configuration := config.Init()

	fmt.Printf("Configuration: %+v\n", configuration)

	eventChan := make(chan event.Event, len(configuration.Cameras))

	mqtt := MQTT.Init(configuration.MQTT, eventChan)


	for i, camConfig := range configuration.Cameras {
		cam, err := camera.Init(camConfig, eventChan)
		if err != nil {
			panic(fmt.Errorf("init: cam %d (%s/[%s] error", i, camConfig.Host, camConfig.Topic))
		}
		cam.ReceiveEvents()
	}

	mqtt.SendEvents()
}
