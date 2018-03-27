package main

import (
	"dahuaevents2mqtt/camera"
	"dahuaevents2mqtt/config"
	"fmt"
)

func main() {

	configuration := config.Init()

	fmt.Printf("Configuration: %+v\n", configuration)

	camChan := make(chan camera.Event, 0)

	cam := camera.Init(camera.Config{
		Topic:    "/openhab/parking_motion/state/set",
		Map:      map[string]string{
			"Start": "ON",
			"Stop": "OFF",
		},
		Host:     "192.168.1.194",
		Port:     "80",
		Username: "admin",
		Password: "admin",
		Events:   []string{"VideoMotion"},
	}, camChan)
	cam.Do()

	for {
		select {
		case event := <- camChan:
			fmt.Printf("%+v\n", event)
		}
	}
}
