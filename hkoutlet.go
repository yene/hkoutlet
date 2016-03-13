package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

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
	o map[string]*accessory.Switch
}

var addr = flag.String("addr", "192.168.1.15:5001", "Pilight daemon IP")

func main() {
	flag.Parse()

	devices.d = make(map[string]Device)
	devices.o = make(map[string]*accessory.Switch)

	u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Println("reconnect here")
	}
	defer ws.Close()

	getConfig(ws)
	initalValues(ws)

	switches := make([]*accessory.Accessory, 0)

	for _, d := range devices.d {
		name := d.Devices[0]
		info := model.Info{
			Name:         name,
			Manufacturer: "Intertechno",
			Model:        "IT-1500",
		}

		sw := accessory.NewSwitch(info)
		sw.SetOn(isOn(d.Values.State))
		sw.OnStateChanged(func(on bool) {
			if on == true {
				turnOn(ws, name)
			} else {
				turnOff(ws, name)
			}
		})

		devices.o[name] = sw
		switches = append(switches, sw.Accessory)
	}
	// Fake accessory to set the device name. Name cannot contain space.
	label := accessory.New(model.Info{Name: "Switch"})
	pin := "00102003"
	t, err := hap.NewIPTransport(hap.Config{Pin: pin}, label, switches...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pin is", pin)

	go listenForUpdates(ws)

	hap.OnTermination(func() {
		t.Stop()
	})

	t.Start()

}

func turnOn(ws *websocket.Conn, name string) {
	err := ws.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"control\",\"code\":{\"device\":\""+name+"\",\"state\":\"on\"}}"))
	if err != nil {
		log.Fatal(err)
	}
}

func turnOff(ws *websocket.Conn, name string) {
	err := ws.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"control\",\"code\":{\"device\":\""+name+"\",\"state\":\"off\"}}"))
	if err != nil {
		log.Fatal(err)
	}
}

func listenForUpdates(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error::: %s", err.Error())
			continue
		}

		if strings.Contains(string(message), "\"origin\":\"update\"") {
			var d Device
			if err := json.Unmarshal(message, &d); err != nil {
				log.Fatal(err)
			}

			name := d.Devices[0]
			devices.Lock()
			devices.d[name] = d
			devices.o[name].SetOn(isOn(d.Values.State))
			devices.Unlock()
		}

	}
}

func getConfig(ws *websocket.Conn) {
	err := ws.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"request config\"}"))
	if err != nil {
		log.Fatal(err)
	}
	_, message, err := ws.ReadMessage()
	err = ioutil.WriteFile("config.json", message, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func initalValues(ws *websocket.Conn) {
	err := ws.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"request values\"}"))
	if err != nil {
		log.Fatal(err)
	}
	_, message, err := ws.ReadMessage()
	var ds []Device
	if err := json.Unmarshal(message, &ds); err != nil {
		log.Fatal(err)
	}

	for _, d := range ds {
		name := d.Devices[0]
		devices.d[name] = d
	}

}

func debug() {
	for _, d := range devices.d {
		log.Printf("Device: %s %s\n", d.Devices[0], d.Values.State)
	}
}

func isOn(s string) bool {
	if s == "on" {
		return true
	} else {
		return false
	}
}
