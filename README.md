# hkoutlet
Turns [pilight](https://www.pilight.org/) outlets into HomeKit accessory using 
[HomeControl](https://github.com/brutella/hc).

I am using Intertechno IT-1500 outlets together with a 433Mhz sender on the Raspberry Pi.

![pi with 433](pi.jpg)


# TODO
* pretty pin dialog like homebridge, that works with cam
* use built in JSON methods https://github.com/golang-samples/websocket/blob/master/websocket-chat/src/chat/client.go#L101
* Remove unneeded map

# License
CC BY Yannick Weiss

# Pilight API
## Get Config
`{"action":"request config"}`

## Get Inital Values
`{"action":"request values"}`

## Updates
`{"origin":"update","type":1,"devices":["Switch1"],"values":{"timestamp":1456200104,"state":"off"}}`

## Change Value
`{"action":"control","code":{"device":"Switch1","state":"on"}}`

## Questions
What does type 1 mean in the update?


# Credits
* https://github.com/knalli/homebridge-pilight
* https://github.com/brutella/hc
* https://www.pilight.org/
