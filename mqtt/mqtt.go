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

package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"dahuaevents2mqtt/event"
	"time"
	"github.com/jvehent/service-go"
)

type Config struct {
	Host string
	Port string
	Timeout int
}

type mqtt struct {
	config Config
	opts *MQTT.ClientOptions
	client MQTT.Client
	eventChan chan event.Event
	log service.Logger
}

func Init(config Config, eventChan chan event.Event, log service.Logger) *mqtt {
	mqtt := new(mqtt)
	mqtt.config = config
	mqtt.opts = MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", mqtt.config.Host, mqtt.config.Port))
	mqtt.eventChan = eventChan
	mqtt.log = log
	return mqtt
}

func (mqtt mqtt) SendEvents(exit chan struct{}) {
	f := func() {
		log := mqtt.log
		connected := false
		for {
			if !connected {
				mqtt.client = MQTT.NewClient(mqtt.opts)
				if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
					log.Error("connect error: %+v", token.Error())
					time.Sleep(1 * time.Second)
					log.Error("Reconnect...")
					continue
				} else {
					connected = true
				}
			}
			select {
			case e := <-mqtt.eventChan:
				log.Info("Event to publish: %+v", e)
				if token := mqtt.client.Publish(e.Topic, 0, false, e.Action); token.Wait() && token.Error() != nil {
					log.Error("Publish error: %+v", token.Error())
					time.Sleep(1 * time.Second)
					continue
				}
			case <- exit:
				return
			}
		}
	}
	f()
}