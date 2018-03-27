package camera

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Topic    string
	Map      map[string]string
	Host     string
	Port     string
	Username string
	Password string
	Events   []string
}

type camera struct {
	config       Config
	client       *http.Client
	connected    bool
	eventChannel chan Event
}

type Event struct {
	Topic  string
	Camera string
	Code   string
	Action string
	Index  string
}

func Init(config Config, channel chan Event) *camera {
	camera := new(camera)
	camera.config = config
	camera.client = &http.Client{}
	camera.eventChannel = channel
	return camera
}

func (camera *camera) Do() {

	const dataBlockSize = 1024

	events := strings.Join(camera.config.Events, ",")

	url := fmt.Sprintf("http://%s:%s@%s:%s/cgi-bin/eventManager.cgi?action=attach&codes=[%s]",
		camera.config.Username,
		camera.config.Password,
		camera.config.Host,
		camera.config.Port,
		events)
	log.Printf("URL: %s", url)
	go func() {
		for {
			var res *http.Response
			if !camera.connected {
				var err error
				if res, err = camera.client.Get(url); err != nil {
					log.Printf("Get error: %+v", err)
					time.Sleep(1 * time.Second)
					log.Printf("Reconnect...")
					continue
				}
			}
			log.Printf("Header: %+v\n", res.Header)
			log.Println("Read data")
			result := make([]byte, 0)
			data := make([]byte, dataBlockSize)
			for {
				n, err := res.Body.Read(data)
				if err != nil {
					log.Printf("Read error: %+v", err)
					camera.connected = false
					log.Printf("Reconnect...")
					break
				}
				if n > 0 {
					result = append(result, data[:n]...)
				}
				if n != dataBlockSize {
					break
				}
			}
			log.Printf("Body:\n%+v", string(result))
			event := Event{
				Camera: camera.config.Host,
				Topic: camera.config.Topic,
			}
			for _, s := range strings.Split(string(result), "\r\n") {
				if strings.HasPrefix(s, "Code=") {
					for _, s1 := range strings.Split(s, ";") {
						r := strings.Split(s1, "=")
						switch r[0] {
						case "Code":
							event.Code = r[1]
						case "action":
							event.Action = r[1]
							event.Action = camera.config.Map[event.Action]
						case "index":
							event.Index = r[1]
						}
					}
					camera.eventChannel <- event
				}
			}
		}
	}()
}
