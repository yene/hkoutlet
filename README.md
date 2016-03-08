# hkoutlet

[HomeControl](https://github.com/brutella/hc)

[Outlet](https://github.com/brutella/hc/blob/master/model/accessory/outlet.go)

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
