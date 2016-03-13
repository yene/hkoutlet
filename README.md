# hkoutlet
Turns [pilight](https://www.pilight.org/) outlets into HomeKit accessory using 
[HomeControl](https://github.com/brutella/hc).

I am using Intertechno IT-1500 outlets together with a 433Mhz sender on the Raspberry Pi.

![pi with 433](pi.jpg)

## Bug NewIPTransport Bridge name
"The first accessory acts as the HomeKit bridge"

If the name contains space it does not work. It works in HomeBridge.

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
* pilight: What does type 1 mean in the update?
* How does a websocket require credentials?
* Can pilight give the devices names which siri can use.

## Credits
* https://github.com/knalli/homebridge-pilight
* https://github.com/brutella/hc
* https://www.pilight.org/
