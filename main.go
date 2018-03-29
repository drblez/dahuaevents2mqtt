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
