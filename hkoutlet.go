package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

var devices struct {
	sync.Mutex
	d map[string]Device
	o map[string]*accessory.Switch
}

var config PilightConfig

var addr = flag.String("addr", "192.168.1.15:5001", "Pilight daemon IP")
var pin = flag.String("pin", "00102003", "HomeKit pin (8 digits)")
var bridgeName = flag.String("name", "Switch", "Bridgename (no space)")

func main() {
	flag.Parse()

	devices.d = make(map[string]Device)
	devices.o = make(map[string]*accessory.Switch)

	u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	for err != nil {
		log.Println("error:", err.Error())
		time.Sleep(time.Second * 3)
		ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	}

	defer ws.Close()

	getConfig(ws)
	initalValues(ws)

	switches := make([]*accessory.Accessory, 0)

	for _, d := range devices.d {
		id := d.Devices[0]
		name := d.Devices[0]
		// Use the GUI name if defined.
		if val, ok := config.Gui[id]; ok {
			name = val.Name
		}
		info := model.Info{
			Name:         name,
			Manufacturer: "Intertechno",
			Model:        "IT-1500",
		}

		sw := accessory.NewSwitch(info)
		sw.SetOn(isOn(d.Values.State))
		sw.OnStateChanged(func(on bool) {
			if on == true {
				turnOn(ws, id)
			} else {
				turnOff(ws, id)
			}
		})

		devices.o[id] = sw
		switches = append(switches, sw.Accessory)
	}
	// Fake accessory to set the device name. Name cannot contain space.
	label := accessory.New(model.Info{Name: *bridgeName}, accessory.TypeSwitch)
	t, err := hap.NewIPTransport(hap.Config{Pin: *pin}, label, switches...)
	if err != nil {
		log.Fatal(err)
	}
	printPin(*pin)

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
		if err != nil { // TODO: this error should be checekd for disconnect.
			for err != nil {
				log.Println("error:", err.Error())
				time.Sleep(time.Second * 3)
				u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
				ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
			}
			continue
		}

		if strings.Contains(string(message), "\"origin\":\"update\"") {
			var d Device
			if err := json.Unmarshal(message, &d); err != nil {
				log.Fatal(err)
			}

			id := d.Devices[0]
			devices.Lock()
			devices.d[id] = d
			devices.o[id].SetOn(isOn(d.Values.State))
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
	//err = ioutil.WriteFile("config.json", message, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(message, &config); err != nil {
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
		id := d.Devices[0]
		devices.d[id] = d
	}
}

func printPin(pin string) {
	fmt.Println("Scan this code with your HomeKit App on your iOS device to pair with Homebridge:")

	pin = pin[0:3] + "-" + pin[3:5] + "-" + pin[5:]
	fmt.Printf("\x1b[30;47m%s\x1b[0m\n", "                       ")
	fmt.Printf("\x1b[30;47m%s\x1b[0m\n", "    ┌────────────┐     ")
	fmt.Printf("\x1b[30;47m%s\x1b[0m\n", "    │ "+pin+" │     ")
	fmt.Printf("\x1b[30;47m%s\x1b[0m\n", "    └────────────┘     ")
	fmt.Printf("\x1b[30;47m%s\x1b[0m\n", "                       ")
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

type Device struct {
	Devices []string
	Values  struct {
		Timestamp int
		State     string // "on" or "off"
	}
}

type PilightConfig struct {
	Gui map[string]struct {
		Name string `json:"name"`
	} `json:"gui"`
}
