#FabSmith
Internal machine activation software for Fab Labs. Build with [BeeGo](http://beego.me) framework for [GoLang](https://golang.org) and [Angular.js](https://angularjs.org).

##Table of contents
- [Quick-start](#quick-start)
  - [Configuration](#configuration)
  - [Compiling and Installing Go](#compiling-and-installing-go)
  - [Beego](#beego)
  - [AngularJS](#angularjs)
  - [Node.js](#nodejs)
- [Versioning](#versioning)
- [Runmode](#runmode)
- [Hexabus](#hexabus)
  - [Requirements](#requirements)
  - [Hexabus IPv6 Network Setup](#hexabus-ipv6-network-setup)
- [Development](#development)

##Quick-start
If you have [GoLang](https://golang.org), [Bee](https://github.com/beego/bee) and [Beego](http://beego.me), use the follwing to clone the repo:  
```
go get github.com/kr15h/fabsmith
```
You can find the project in `$GOPATH/src/github.com/kr15h/fabsmith` once cloned.

Make sure you have [Bower](http://bower.io), run `bower install` from within the `static/` directory to install all the front-end dependencies.

Current version of this software is being tested on a Raspberry Pi, Raspbian and this README should be compatible with other Linux systems.

For Hexabus part now there is working solution that requires the use of `radvd`. The next step would be to understand how to make it work on Mac OS X with `rtadvd`. 

### Configuration
You have to make a copy of the `conf/app.example.conf` file to be able to run Fabsmith:  
```
cp conf/app.example.conf conf/app.conf
```

Compile FabSmith with `bee run` from the `github.com/kr15h/fabsmith` direcotry. It will compile and run the project. To use just the binary afterwards, use the following:  
```
sudo ./fabsmith
```

You can set the runmode via environment variables `FABSMITH_RUNMODE` or `BEEGO_RUNMODE` to alter the runmode. You can pass environment variables directly on launching the program:  
```
sudo FABSMITH_RUNMODE="prod" ./fabsmith
```

The `FABSMITH_RUNMODE` environment variable overrides the config runmode. The default runmode is `dev`.

###Compiling and Installing Go
You will need to compile GoLang from source on the Raspberry Pi. Takes about 2 hours. 

 1. Memory split. You should give more for the CPU of the Pi. Open `sudo raspi-config`, go to **Advanced Settings**, select **Memory Split** and enter **128**. On a 512M Raspberry Pi 128M will be given to the GPU and the rest to CPU. This should be enough.
 
 2. Swap space. Do this if you have less than 512MB of RAM. Raspberry Pi Model B or B+ is recomended - it has 512 MB of RAM and if you assign only 64MB to the GPU, the rest should be enough to compile Go. Do the following to create some swap space:  
 ```
 % sudo dd if=/dev/zero of=/import/nas/swap bs=1024 count=1048576
1048576+0 records in
1048576+0 records out
1073741824 bytes (1.1 GB) copied, 136.045 s, 7.9 MB/s
% sudo losetup /dev/loop0 /import/nas/swap
% sudo mkswap /dev/loop0
Setting up swapspace version 1, size = 1048572 KiB
no label, UUID=7ba9443d-c64c-416f-9931-39e3e2decf0f
% sudo swapon /dev/loop0
% free -m
             total used free shared buffers cached
Mem:           232   78  153      0       0     24
-/+ buffers/cache:   52  179
Swap:         1123   15 1108
 ```  
 Create `/import/nas/swap` if needed with `touch /import/nas/swap`. Use sudo if necessary.
 
 3. Pre-requisites  
 ```
 % sudo apt-get install -y mercurial gcc libc6-dev
 ```
 
 4. Clone source  
 ```
 % hg clone -u default https://code.google.com/p/go $HOME/go
warning: code.google.com certificate with fingerprint 9f:af:b9:ce:b5:10:97:c0:5d:16:90:11:63:78:fa:2f:37:f4:96:79 not verified (check hostfingerprints or web.cacerts config setting)
destination directory: go
requesting all changes
adding changesets
adding manifests
adding file changes
added 14430 changesets with 52478 changes to 7406 files (+5 heads)
updating to branch default
3520 files updated, 0 files merged, 0 files removed, 0 files unresolved
 ```
 
 5. Build
 ```
 % cd $HOME/go/src
% ./all.bash
 ```
 When done, move `~/go` to `/opt/go`. This is where user programs that are more or less self-contained should go.
 
 6. Setup. Open your shell config file (`.bashrc` or `.zshrc` - depends on what are you using) and add the following:  
 ```
 export GOROOT=/opt/go
export GOPATH=$HOME/go-workspace
export PATH="/usr/local/bin:/usr/bin:/bin:/usr/bin/X11:/usr/games:$GOROOT/bin:$GOPATH/bin"
 ```
 Create `~/go-workspace` directory. Name it as you wish, but remember to change the path to it in the lines above added to the shell config file.
 
These instructions have been adapted from [Dave Cheney's](http://dave.cheney.net/2012/09/25/installing-go-on-the-raspberry-pi) blog article. If something goes wrong - refer to it.
 

###Beego
You will need to install BeeGo MVC framework to make use of the code found in this repository. Refer to the official [BeeGo installation instuctions](http://beego.me/quickstart) to do that.

###AngularJS
On the Beego side there is single frontend template `views/index.html` that loads in AngularJS application from within `static/` directory. The AngularJS side is build on the [AngularJS Seed](https://github.com/angular/angular-seed) project.

Once you clone the source from this repository, you have to cd into the `static/` directory and execute `npm install` to install all the project spceific dependencies.

###Node.js
We use Node.js to fully benefit from the AngularJS Seed project. Use the [Node Version Manager](https://github.com/creationix/nvm) to install latest Node.js version. On the Raspberry Pi it will compile it from source and it takes approximately 2 hours.

###MySQL

Install on Linux:  
```
sudo apt-get update
sudo apt-get install mysql-server mysql-client
```

Create a database with the `mysql` tool (`mysql -u root -p`):
```
CREATE DATABASE fabsmith DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
```

Create a safe MySQL user (replace username and password): 
```
GRANT ALL PRIVILEGES ON fabsmith.* To 'fabsmith'@'localhost' IDENTIFIED BY 'fabsmith';
SET PASSWORD FOR 'fabsmith'@'localhost' =  PASSWORD('fabsmith');
```

Restore a database from a dump:
```
./scripts/restoredb.sh
```

In future the Beego application should take care of this on it's own.

Dump database:
```
./scripts/dumpdb.sh
```

Edit the files `restoredb.sh` and `dumpdb.sh` to add your specific username and password combination.

Create a safe MySQL user: 
```
GRANT ALL PRIVILEGES ON fabsmith.* To 'fabsmith'@'localhost' IDENTIFIED BY 'fabsmith';
```

##Versioning
FabSmith will benefit from semantic versioning. Read about it [here](http://semver.org).

##Hexabus

We use [Hexabus](https://github.com/mysmartgrid/hexabus) GoLang library to communicate with the remote switches from within the Beego application.

###Requirements
- Go 1.3+
- radvd
- [This Guide](https://github.com/mysmartgrid/hexabus/wiki/Connect-PC-Directly)

It might be necessary to sync devices.  
1. Press button on USB router while it starts blinking slowy  
2. Press button on the Switch device and wait till the light becomes green


###Hexabus IPv6 Network Setup

1. You need the USB stick that comes tohether with the Hexabus switches. Write down the IPv6 address on it - it should look like the following:
```
fe80::50:c4ff:fe04:01B8
```

2. Replace the `fe80` part with your own **hex** prefix:
```
fefe::50:c4ff:fe04:01B8
```  

3. Type `sudo ifconfig usb0 add fefe::50:c4ff:fe04:01B8/64` and add the following to `/etc/network/interfaces`:  
```
auto usb0
allow-hotplug usb0
iface usb0 inet6 static
    ip6addr fefe::50:c4ff:fe04:01B8
    netmask 64
```

4. Type `sudo ifconfig` to see if an IPv6 IP has been applied. You should see output similar to this:
```
usb0      Link encap:Ethernet  HWaddr 02:50:c4:04:02:bd
      inet6 addr: fe80::50:c4ff:fe04:01B8/64 Scope:Link
      UP BROADCAST RUNNING MULTICAST  MTU:1284  Metric:1
      RX packets:10 errors:0 dropped:0 overruns:0 frame:0
      TX packets:8 errors:0 dropped:0 overruns:0 carrier:0
      collisions:0 txqueuelen:1000
      RX bytes:460 (460.0 B)  TX bytes:648 (648.0 B)
```

5. If you do not have the `inet6 addr` line, probably it's because the `ipv6` kernel module is not enabled. To enable it, type `sudo modprobe ipv6` and add `ipv6` to `/etc/modules` so it is being enabled after the next reboot.  

6. Pair the USB plug with the switch. Press and hold the button on the plug until the LED starts to blink red. Then press and hold the button on the power switch for a couple of seconds - it will blink red and then green on release. The LED on the USB plug should turn green as well. If both LEDs (the one on the USB switch and on the power plug) are green - the devices have been paired.

7. Try to ping the plug:  
```
ping6 -Iusb0 fe80::50:c4ff:fe04:8390
```  
Use the original plug address here.

8. Install `radvd` with `sudo apt-get install radvd`. It stands for Router Advertisement Daemon. Note that you need only one `radvd` running on your IPv6 network.

9. Config `radvd` **(if it is the only instance on the IPv6 network)** by opening the file `/etc/radvd.conf`. Add the following contents there:  
```
interface usb0
{
AdvSendAdvert on;
AdvLinkMTU 1280;
AdvCurHopLimit 128;
AdvReachableTime 360000;
MinRtrAdvInterval 100;
MaxRtrAdvInterval 150;
AdvDefaultLifetime 200;
prefix fafa::/64
{
  AdvOnLink on;
  AdvAutonomous on;
  AdvPreferredLifetime 4294967295;
  AdvValidLifetime 4294967295;
};
};
```

10. Ping your plugs with global addresses:  
```
ping6 fafa::50:c4ff:fe04:8390
```  
As you see we don't use the `-Iusb0` part anymore.

11. We are set at this point - associate the Hexabus device IPv6 addresses with the machines in the FabSmith database.


## Development

The development environment is a constant work in progress and does not implement a decent test-driven development workflow yet.  

The Angular JS applications that can be found in the `views/` directory are integrated with a custom [Grunt](http://gruntjs.com) workflow. While developing use the `grunt dev` task - it will constantly check for changes in your JavaScript and [Less](http://lesscss.org) files, compile CSS from Less and check JavaScript for erros with the help of [JSHint](http://jshint.com). You can use a [LiveReload](http://feedback.livereload.com/knowledgebase/articles/86242-how-do-i-install-and-use-the-browser-extensions-) browser extension to reload views each time a change is detected.

When ready to move the system to `prod` runmode, run `grunt prod` to compile the production mode of the Angular JS applications.

More info about the development workflow will be added to Wiki.