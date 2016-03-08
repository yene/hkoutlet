package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

var origin = "http://localhost/"
var url = "ws://192.168.100.142:5001/"

type Device struct {
	Devices []string
	Values  struct {
		Timestamp int
		State     string // "on" or "off"
	}
}

var devices struct {
	sync.Mutex
	d map[string]Device
}

func turnOn(ws *websocket.Conn) {
	log.Println("Turn Light On")
	err := websocket.Message.Send(ws, "{\"action\":\"control\",\"code\":{\"device\":\"Switch1\",\"state\":\"on\"}}")
	if err != nil {
		log.Fatal(err)
	}
}

func turnOff(ws *websocket.Conn) {
	log.Println("Turn Light Off")
	err := websocket.Message.Send(ws, "{\"action\":\"control\",\"code\":{\"device\":\"Switch1\",\"state\":\"off\"}}")
	if err != nil {
		log.Fatal(err)
	}
}

func listenForUpdates(ws *websocket.Conn) {
	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}

		if strings.Contains(message, "\"origin\":\"update\"") {
			var d Device
			if err := json.Unmarshal([]byte(message), &d); err != nil {
				log.Fatal(err)
			}

			devices.Lock()
			defer devices.Unlock()
			name := d.Devices[0]
			// We have to compare here for change I think
			devices.d[name] = d
		}
	}
}

func getConfig(ws *websocket.Conn) {
	err := websocket.Message.Send(ws, "{\"action\":\"request config\"}")
	if err != nil {
		log.Fatal(err)
	}
	var answer []byte
	websocket.Message.Receive(ws, &answer)
	err = ioutil.WriteFile("config.json", answer, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func initalValues(ws *websocket.Conn) {
	err := websocket.Message.Send(ws, "{\"action\":\"request values\"}")
	if err != nil {
		log.Fatal(err)
	}
	var answer []byte
	websocket.Message.Receive(ws, &answer)
	var ds []Device
	if err := json.Unmarshal(answer, &ds); err != nil {
		log.Fatal(err)
	}

	for _, d := range ds {
		name := d.Devices[0]
		devices.d[name] = d
	}
	debug()
}

func debug() {
	for _, d := range devices.d {
		fmt.Printf("Device: %s %s\n", d.Devices[0], d.Values.State)
	}
}

func main() {
	devices.d = make(map[string]Device)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to ", url)

	getConfig(ws)
	initalValues(ws)

	info := model.Info{
		Name:         "Outlet",
		Manufacturer: "Intertechno",
		Model:        "IT-1500",
	}

	// https://github.com/brutella/hc/blob/master/model/accessory/outlet.go
	outlet := accessory.NewOutlet(info)
	// TODO: set inital state

	outlet.OnStateChanged(func(on bool) {
		if on == true {
			turnOn(ws)
		} else {
			turnOff(ws)
		}
	})

	pin := "00102003"
	t, err := hap.NewIPTransport(hap.Config{Pin: pin}, outlet.Accessory)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pin is ", pin)

	go listenForUpdates(ws)

	hap.OnTermination(func() {
		t.Stop()
	})

	t.Start()

}
