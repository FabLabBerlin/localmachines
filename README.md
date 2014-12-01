#FabSmith
Internal machine activation software for Fab Labs. Build with [BeeGo](http://beego.me) framework for [GoLang](https://golang.org) and [Angular.js](https://angularjs.org).

##Table of contents
- [Quick start](#quick-start)
- [Versioning](#versioning)

##Quick start
This solution is supposed to be run on a computer like the Raspberry Pi to benefit from it's GPIO pins - it's where the external hardware has to be connected in order to get full support of FabSmith's features.

You will need to install BeeGo MVC framework to make use of the code found in this repository. Refer to the official [BeeGo installation instuctions](http://beego.me/quickstart) to do that.

##Versioning
FabSmith will benefit form semantic versioning. Read about it [here](http://semver.org).

##Hexabus requirements

- Go 1.3+
- radvd
- [This Guide](https://github.com/mysmartgrid/hexabus/wiki/Connect-PC-Directly)

It might be necessary to sync devices.  
1. Press button on USB router while it starts blinking slowy  
2. Press button on the Switch device and wait till the light becomes green



