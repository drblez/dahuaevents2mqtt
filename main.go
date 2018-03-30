/*
 * Copyright 2018 Dr. Blez AKA Ruslan Stepanenko
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"dahuaevents2mqtt/event"
	"dahuaevents2mqtt/camera"
	"dahuaevents2mqtt/config"
	MQTT "dahuaevents2mqtt/mqtt"
	"github.com/jvehent/service-go"
	"os"
)

var log service.Logger
var exit = make(chan struct{})

func main() {

	var name = "dahuaevents2mqtt"
	var displayName = "dahuaevents2mqtt"
	var desc = "Send Dahua IPC events to MQTT broker"

	var s, err = service.NewService(name, displayName, desc)
	log = s

	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		fmt.Printf("Verb is %s\n", verb)
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			do(exit, log, true)
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		fmt.Println("Start service")
		go do(exit, log, false)
		return nil
	}, func() error {
		// stop
		exit <- struct{}{}
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}

}

func do(exit chan struct{}, log service.Logger, runCommand bool) {

	configuration := config.Init()

	fmt.Printf("Configuration: %+v\n", configuration)

	eventChan := make(chan event.Event, len(configuration.Cameras))

	mqtt := MQTT.Init(configuration.MQTT, eventChan, log)


	for i, camConfig := range configuration.Cameras {
		cam, err := camera.Init(camConfig, eventChan, log)
		if err != nil {
			panic(fmt.Errorf("init: cam %d (%s/[%s] error", i, camConfig.Host, camConfig.Topic))
		}
		cam.ReceiveEvents()
	}

	mqtt.SendEvents(exit)

	if runCommand {
		for {}
	}
}