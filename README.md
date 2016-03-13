# hkoutlet
Turns [pilight](https://www.pilight.org/) outlets into HomeKit accessory using 
[HomeControl](https://github.com/brutella/hc).

I am using Intertechno IT-1500 outlets together with a 433Mhz sender on the Raspberry Pi.

![pi with 433](pi.jpg)

## TODO
* change websocket code to gorilla websocket
* change outlet to switch
* test with siri
* websocket should reconnect, always, test


----
* pretty pin dialog like homebridge, that works with iOS cam
* use built in JSON methods https://github.com/golang-samples/websocket/blob/master/websocket-chat/src/chat/client.go#L101
* Remove unneeded map
* Integrate the pilight project, the part I need directly here.
* when refused keep retrying to reconnect but inform the user
* when connection refused don't fail, reconnect
* compare project with homebridge pilight

## License
CC BY Yannick Weiss

## Pilight API
### Get Config
`{"action":"request config"}`

### Get Inital Values
`{"action":"request values"}`

### Updates
`{"origin":"update","type":1,"devices":["Switch1"],"values":{"timestamp":1456200104,"state":"off"}}`

### Change Value
`{"action":"control","code":{"device":"Switch1","state":"on"}}`

### Open Questions
What does type 1 mean in the update?

## Credits
* https://github.com/knalli/homebridge-pilight
* https://github.com/brutella/hc
* https://www.pilight.org/
