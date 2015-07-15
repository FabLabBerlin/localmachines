# FabSmith
Internal machine activation software for Fab Labs. Build with [BeeGo](http://beego.me) framework for [GoLang](https://golang.org) and [Angular.js](https://angularjs.org).

## Table of contents
- [Quick-start](#quick-start)
  - [Configuration](#configuration)
  - [Running](#running)
  - [Compiling and Installing Go](#compiling-and-installing-go)
  - [Installing Beego](#installing-beego)
  - [Node.js](#nodejs)
- [Versioning](#versioning)
- [Runmode](#runmode)
- [Hexabus](#hexabus)
  - [Requirements](#requirements)
  - [Hexabus IPv6 Network Setup](#hexabus-ipv6-network-setup)
- [Development](#development)
  - [Core Development](#core-development)
  - [Clients Development](#clients-development)
- [API Documentation](#api-documentation)

## Quick-Start
This section will show you how to get it working as fast as possible.

If you have [GoLang](https://golang.org), [Bee](https://github.com/beego/bee) and [Beego](http://beego.me), use the follwing to clone the repo:  
```
go get github.com/kr15h/fabsmith
```

You can find the project in `$GOPATH/src/github.com/kr15h/fabsmith` once cloned. Navigate to it using `cd $GOPATH/src/github.com/kr15h/fabsmith`.

Run `go get` to install all GoLang dependencies.

### Configuration
You have to make a copy of the `conf/app.example.conf` file to be able to run Fabsmith:  
```
cp conf/app.example.conf conf/app.conf
```

Edit the `conf/app.conf` file to match your environment.

Compile FabSmith with `bee run` from the `github.com/kr15h/fabsmith` direcotry. It will compile and run the project. To use just the binary afterwards, use the following:  
```
sudo ./fabsmith
```

It might be possible that it complains about the GOPATH not being set. In that case run the binary like this:
```
sudo GOPATH="/home/youruser/gospace" ./fabsmith
```

You can set the runmode via environment variable or `BEEGO_RUNMODE` to alter the runmode. You can pass environment variables directly on launching the program:  
```
sudo BEEGO_RUNMODE="prod" ./fabsmith
```

The `BEEGO_RUNMODE` environment variable overrides the config runmode.

### Running
To run compile and run use `bee run`. It should spawn a local web server accessible through the port you defined in `conf/app.conf`, e.g. `http://localhost:8080`. You can access the admin interface via `http://localhost:8080/admin`.

If you are not able to compile and run at this point - check your config file and whether you are not missing GoLang or Beego.

### Compiling and Installing Go
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

###Installing Beego
You will need to install BeeGo MVC framework to make use of the code found in this repository. Refer to the official [BeeGo installation instuctions](http://beego.me/quickstart) to do that.

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
mysql -u user -p -d database < fabsmith.sql
```

In future the Beego application should take care of this on it's own.

Dump database:
```
mysqldump -u user -p fabsmith > fabsmith.sql
```

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

The development environment is a constant work in progress and does not implement a decent test-driven development workflow yet. Soon.

When ready to move the system to `prod` runmode, run `grunt prod` to compile the production mode of the Angular JS applications.

More info about the development workflow will be added to Wiki.

### Core Development

Core / API development. Everything that is no in the `clients` directory is relevant to this part.

### Clients Development

There is a clean separation between the Fabsmith core engine written in GoLang and the clients that operate with the REST API of the core. Clients are HTML-based and reside under the `/clients` directory.

To switch to development mode, you have to change the `runmode` in `conf/app.conf` to `dev`. Each of the clients have a `dev` and `prod` directory with the same hierarchy. When the runmode is `dev`, clients are loaded from their `dev` directory. If the runmode is `prod`, clients are loaded from their `prod` directories. The `prod` directories of the clients contain minified and optimized files for the client interfaces.

Each of the clients is a separate Angular JS application with a [Grunt](http://gruntjs.com)-based workflow. You need to run the following commands from within each client directory separately to download all the development dependencies and be able to use the clients in `dev` mode:

```
cd clients/machines
npm instal && bower install
cd ../admin
npm instal && bower install
```

Make sure that you have Bower installed (it has been added as one of the npm dependencies though) - if not, follow the instructions on [Bower](http://bower.io) website.

Whenever you are finished with the client development, remember to run `grunt prod` for each of the clients to compile the `prod` version of the client. Also remember to edit the `Gruntfile.js` before running grunt to add extra configuration in case extra libraries have been added and new modules have to be copied.

## API Documentation

We use [Automated API Document](http://beego.me/docs/advantage/docs.md) feature of the Beego framework. It makes use of Beego router namespaces and documentation comments in the router and controller files.

To view the documetntation, compile and run with the following command:

```
bee run watchall true -downdoc=true -gendoc=true
```

It seems that it works also with plain:

```
bee run
```

Then access API documentation via:

```
http://localhost:8080/swagger
```

The port number and host is the same you have set in your `config/app.conf` file.

After you have done some changes to the code, run the following to update API documentation:

```
bee generate docs
```
