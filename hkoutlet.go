package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

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
	o map[string]*accessory.Outlet
}

func turnOn(ws *websocket.Conn, name string) {
	//log.Println("Turn Light On", name)
	err := websocket.Message.Send(ws, "{\"action\":\"control\",\"code\":{\"device\":\""+name+"\",\"state\":\"on\"}}")
	if err != nil {
		log.Fatal(err)
	}
}

func turnOff(ws *websocket.Conn, name string) {
	//log.Println("Turn Light Off", name)
	err := websocket.Message.Send(ws, "{\"action\":\"control\",\"code\":{\"device\":\""+name+"\",\"state\":\"off\"}}")
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
			devices.d[name] = d
			devices.o[name].SetOn(isOn(d.Values.State))
		}
	}
	fmt.Println("I did exit")
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
}

func debug() {
	for _, d := range devices.d {
		fmt.Printf("Device: %s %s\n", d.Devices[0], d.Values.State)
	}
}

func main() {
	devices.d = make(map[string]Device)
	devices.o = make(map[string]*accessory.Outlet)

	// TODO: add and test reconnect
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to ", url)

	getConfig(ws)
	initalValues(ws)

	outlets := make([]*accessory.Accessory, 0)

	for _, d := range devices.d {
		name := d.Devices[0]
		info := model.Info{
			Name:         name,
			Manufacturer: "Intertechno",
			Model:        "IT-1500",
		}

		outlet := accessory.NewOutlet(info)
		outlet.SetOn(isOn(d.Values.State))
		outlet.OnStateChanged(func(on bool) {
			if on == true {
				turnOn(ws, name)
			} else {
				turnOff(ws, name)
			}
		})

		devices.o[name] = outlet
		outlets = append(outlets, outlet.Accessory)
	}
	label := accessory.New(model.Info{Name: "Outlet2"})
	pin := "00102003"
	t, err := hap.NewIPTransport(hap.Config{Pin: pin}, label, outlets...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pin is ", pin)

	go listenForUpdates(ws)

	hap.OnTermination(func() {
		t.Stop()
	})

	t.Start()

	//go forever()
	//select {} // block forever

}

func isOn(s string) bool {
	if s == "on" {
		return true
	} else {
		return false
	}
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}
