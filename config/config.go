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
	fileName := os.Getenv("DAHUAEVENTS2MQTT_CONFIG")
	if fileName == "" {
		fileName = configFile
	}
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