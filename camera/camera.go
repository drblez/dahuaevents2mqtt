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

package camera

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"errors"
	"dahuaevents2mqtt/event"
	"github.com/jvehent/service-go"
)

type Config struct {
	Topic    string
	Map      map[string]string
	Host     string
	Port     string
	Username string
	Password string
	Events   []string
	ReconnectTimeout int
}

type camera struct {
	config           Config
	client           *http.Client
	connected        bool
	eventChannel     chan event.Event
	reconnectTimeout time.Duration
	log service.Logger
}

func Init(config Config, channel chan event.Event, log service.Logger) (*camera, error) {
	camera := new(camera)
	camera.config = config
	if camera.config.Topic == "" {
		return nil, errors.New("config: Topic can't be empty")
	}
	if camera.config.Host == "" {
		return nil, errors.New("config: Host can't be empty")
	}
	if camera.config.Events == nil {
		camera.config.Events = []string{"VideoMotion"}
	}
	if camera.config.Port == "" {
		camera.config.Port = "80"
	}
	if camera.config.Username == "" {
		camera.config.Username = "admin"
	}
	if camera.config.Password == "" {
		camera.config.Password = "admin"
	}
	camera.client = &http.Client{}
	camera.client.Transport = &http.Transport{}
	camera.eventChannel = channel
	if config.ReconnectTimeout == 0 {
		camera.reconnectTimeout = 1 * time.Second
	} else {
		camera.reconnectTimeout = time.Duration(config.ReconnectTimeout) * time.Second
	}
	camera.log = log
	return camera, nil
}

func (camera *camera) Connected() bool {
	return camera.connected
}

func (camera *camera) ReceiveEvents() {

	log := camera.log

	const dataBlockSize = 1024

	events := strings.Join(camera.config.Events, ",")

	url := fmt.Sprintf("http://%s:%s@%s:%s/cgi-bin/eventManager.cgi?action=attach&codes=[%s]",
		camera.config.Username,
		camera.config.Password,
		camera.config.Host,
		camera.config.Port,
		events)
	log.Info("%s URL: %s", camera.config.Host, url)
	go func() {
		var res *http.Response
		MAIN:
		for {
			if !camera.connected {
				var err error
				if res, err = camera.client.Get(url); err != nil {
					log.Error("%s Get error: %+v", camera.config.Host, err)
					time.Sleep(1 * time.Second)
					log.Error("%s Reconnect...", camera.config.Host)
					continue MAIN
				} else {
					log.Info("%s Header: %+v", camera.config.Host, res.Header)
					camera.connected = true
					log.Info("%s Connected", camera.config.Host)
				}
			}
			log.Info("%s Read data", camera.config.Host)
			result := make([]byte, 0)
			data := make([]byte, dataBlockSize)
			for {
				n, err := res.Body.Read(data)
				if err != nil {
					log.Error("%s Read error: %+v", camera.config.Host, err)
					camera.connected = false
					log.Error("%s Reconnect...", camera.config.Host)
					continue MAIN
				}
				if n > 0 {
					result = append(result, data[:n]...)
				}
				if n != dataBlockSize {
					time.Sleep(1 * time.Second)
					break
				}
			}
			log.Info("%s Body:\n%+v", camera.config.Host, string(result))
			e := event.Event{
				Camera: camera.config.Host,
				Topic:  camera.config.Topic,
			}
			for _, s := range strings.Split(string(result), "\r\n") {
				if strings.HasPrefix(s, "Code=") {
					for _, s1 := range strings.Split(s, ";") {
						r := strings.Split(s1, "=")
						switch r[0] {
						case "Code":
							e.Code = r[1]
						case "action":
							e.Action = r[1]
							if camera.config.Map != nil {
								e.Action = camera.config.Map[e.Action]
							}
						case "index":
							e.Index = r[1]
						}
					}
					camera.eventChannel <- e
				}
			}
		}
	}()
}
