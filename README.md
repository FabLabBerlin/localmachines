#FabSmith
Internal machine activation software for Fab Labs. Build with [BeeGo](http://beego.me) framework for [GoLang](https://golang.org) and [Angular.js](https://angularjs.org).

##Table of contents
- [Quick start](#quick-start)
  - [Compiling and Installing Go](#compiling-and-installing-go)
  - [Beego](#beego)
  - [AngularJS](#angularjs)
  - [Node.js](#nodejs)
- [Versioning](#versioning)
- [Runmode](#runmode)
- [Hexabus](#hexabus)
  - [Requirements](#requirements)

##Quick start
If you have [GoLang](https://golang.org), [Bee](https://github.com/beego/bee) and [Beego](http://beego.me), use the follwing to clone the repo:  
```
go get github.com/kr15h/fabsmith
```
You can find the project in `$GOPATH/src/github.com/kr15h/fabsmith` once cloned.

Make sure you have [Bower](http://bower.io), run `bower install` to install all the front-end dependencies.

Edit `conf/app.conf` to set your MySQL database settings. Yes, you need a MySQL database up and running as well.

Current version of this software is being tested on a Raspberry Pi, Raspbian and this README should be compatible with other Linux systems.

For Hexabus part now there is working solution that requires the use of `radvd`. The next step would be to understand how to make it work on Mac OS X with `rtadvd`. 

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

Create a database with the `mysql` tool:
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
./restoredb.sh
```

In future the Beego application should take care of this on it's own.

Dump database:
```
./dumpdb.sh
```

Edit the files `restoredb.sh` and `dumpdb.sh` to add your specific username and password combination.

Create a safe MySQL user: 
```
GRANT ALL PRIVILEGES ON fabsmith.* To 'fabsmith'@'localhost' IDENTIFIED BY 'fabsmith';
```

##Runmode
Set the `BEEGO_RUNMODE` environment variable in `/etc/environment`:  
```
BEEGO_RUNMODE=prod
```  

Available options are `dev`, `test` and `prod`. For each of them there is a possiblity to set unique settings in the `conf/app.conf` file.

Fabsmith will not be able to run without being able to get the runmode from the environment variable.

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



