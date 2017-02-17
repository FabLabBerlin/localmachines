# Getting started on the RaspberryPi

## Table of contents
- [Quick-start](#quick-start)
  - [Configuration](#configuration)
  - [Running](#running)
  - [Compiling and Installing Go](#compiling-and-installing-go)
  - [Installing Beego](#installing-beego)
  - [Node.js](#nodejs)

## Quick-Start
This section will show you how to get it working as fast as possible. 

If you have [GoLang](https://golang.org), [Bee](https://github.com/beego/bee) and [Beego](http://beego.me), use the follwing to clone the repo:  
```
go get github.com/FabLabBerlin/localmachines
```

You can find the project in `$GOPATH/src/github.com/FabLabBerlin/localmachines` once cloned. Navigate to it using `cd $GOPATH/src/github.com/FabLabBerlin/localmachines`.

Run `go get` to install all GoLang dependencies.

### Configuration
You have to make a copy of the `conf/app.example.conf` file to be able to run Fabsmith:  
```
cp conf/app.example.conf conf/app.conf
```

Edit the `conf/app.conf` file to match your environment.

Compile FabSmith with `bee run` from the `github.com/FabLabBerlin/localmachines` direcotry. It will compile and run the project. To use just the binary afterwards, use the following:  
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
We use Node.js to install Bower and run the frontend tests.  Node.js LTS is used for reference.

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
mysql -u user -p -d database < fabsmith-dump.sql
```

Dump database:
```
mysqldump -u user -p fabsmith > fabsmith-dump.sql
```

Create a safe MySQL user with limited privileges: 
```
GRANT ALL PRIVILEGES ON fabsmith.* To 'fabsmith'@'localhost' IDENTIFIED BY 'fabsmith';
```
