# hkoutlet

[HomeControl](https://github.com/brutella/hc)

[Outlet](https://github.com/brutella/hc/blob/master/model/accessory/outlet.go)

# TODO
* pretty pin dialog like homebridge, that works with cam
* use built in JSON methods https://github.com/golang-samples/websocket/blob/master/websocket-chat/src/chat/client.go#L101
* Remove unneeded map

# License
CC4-by Yannick Weiss

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
https://github.com/knalli/homebridge-pilight
