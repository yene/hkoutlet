package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

// Inspired by https://github.com/knalli/homebridge-pilight

var origin = "http://localhost/"
var url = "ws://192.168.28.134:5001/"

type Device struct {
	Devices []string
	Values  struct {
		Timestamp int
		State     string
	}
}

type Update struct {
	Devices []string
	Values  struct {
		Timestamp int
		State     string
	}
}

var devices map[string]Device

func turnLightOn() {
	log.Println("Turn Light On")
}

func turnLightOff() {
	log.Println("Turn Light Off")
}

func listenForUpdates(ws *websocket.Conn, updates chan string) {
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

			name := d.Devices[0]
			// We have to compare here for change I think

			devices[name] = d
		}

		//updates <- message
		// Send updates back via channel or aqquire mutex and modify em myself.
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
		devices[name] = d
	}

}

func debug() {
	for _, d := range devices {
		fmt.Printf("Device: %s %s\n", d.Devices[0], d.Values.State)
	}
}

func main() {
	devices = make(map[string]Device)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to ", url)

	getConfig(ws)
	initalValues(ws)
	debug()

	updates := make(chan string)
	go listenForUpdates(ws, updates)

	info := model.Info{
		Name:         "Radio Controlled Outlet",
		Manufacturer: "Intertechno",
	}

	light := accessory.NewLightBulb(info)
	light.OnStateChanged(func(on bool) {
		if on == true {
			turnLightOn()
		} else {
			turnLightOff()
		}
	})

	t, err := hap.NewIPTransport(hap.Config{Pin: "32191123"}, light.Accessory)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pin is 32191123")

	hap.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
