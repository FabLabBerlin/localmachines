#FabSmith
Internal machine activation software for Fab Labs. Build with [BeeGo](http://beego.me) framework for [GoLang](https://golang.org) and [Angular.js](https://angularjs.org).

##Table of contents
- [Quick start](#quick-start)
  - [Beego](#beego)
  - [AngularJS](#angularjs)
  - [Node.js](#nodejs)
- [Versioning](#versioning)
- [Hexabus](#hexabus)
  - [Requirements](#requirements)

##Quick start
This solution is supposed to be run on a computer like the Raspberry Pi to benefit from it's GPIO pins - it's where the external hardware has to be connected in order to get full support of FabSmith's features.

###Beego
You will need to install BeeGo MVC framework to make use of the code found in this repository. Refer to the official [BeeGo installation instuctions](http://beego.me/quickstart) to do that.

###AngularJS
On the Beego side there is single frontend template `views/index.html` that loads in AngularJS application from within `static/` directory. The AngularJS side is build on the [AngularJS Seed](https://github.com/angular/angular-seed) project.

Once you clone the source from this repository, you have to cd into the `static/` directory and execute `npm install` to install all the project spceific dependencies.

###Node.js
We use Node.js to fully benefit from the AngularJS Seed project. Use the [Node Version Manager](https://github.com/creationix/nvm) to install latest Node.js version. On the Raspberry Pi it will compile it from source and it takes approximately 2 hours.

##Versioning
FabSmith will benefit form semantic versioning. Read about it [here](http://semver.org).

##Hexabus

We use [Hexabus](https://github.com/mysmartgrid/hexabus) GoLang library to communicate with the remote switches from within the Beego application.

###Requirements
- Go 1.3+
- radvd
- [This Guide](https://github.com/mysmartgrid/hexabus/wiki/Connect-PC-Directly)

It might be necessary to sync devices.  
1. Press button on USB router while it starts blinking slowy  
2. Press button on the Switch device and wait till the light becomes green



