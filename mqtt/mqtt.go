package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"dahuaevents2mqtt/event"
	"log"
	"time"
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
}

func Init(config Config, eventChan chan event.Event) *mqtt {
	mqtt := new(mqtt)
	mqtt.config = config
	mqtt.opts = MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", mqtt.config.Host, mqtt.config.Port))
	mqtt.eventChan = eventChan
	return mqtt
}

func (mqtt mqtt) SendEvents() {
	connected := false
	for {
		if !connected {
			mqtt.client = MQTT.NewClient(mqtt.opts)
			if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
				log.Printf("connect error: %+v", token.Error())
				time.Sleep(1 * time.Second)
				log.Printf("Reconnect...")
				continue
			} else {
				connected = true
			}
		}
		select {
		case e := <-mqtt.eventChan:
			log.Printf("Event to publish: %+v", e)
			if token := mqtt.client.Publish(e.Topic, 0, false, e.Action); token.Wait() && token.Error() != nil {
				log.Printf("Publish error: %+v", token.Error())
				time.Sleep(1 * time.Second)
				continue
			}
		}

	}
}