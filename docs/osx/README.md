# OS X Development Setup

This section covers how to setup Localmachines on your development machine.

## 1. Install Go

The easiest and quickest way to get started is to download Go binaries from
the [Go Download Page](https://golang.org/dl/). Just download the .pkg file
of the most recent version, 1.5.3 as of this writing, and open it. It will
automatically install itself in `/usr/local/go`.

## 2. Set Go environment variables

Now you also need to setup some environment variables.  E.g. in `$HOME/.profile`:

```
export GOROOT=/usr/local/go
```

More over you need a `GOPATH`.  The `GOPATH` is where all your dependencies go.
Many people choose `$HOME/go` for this. Example:

```
export GOPATH=$HOME/go
```

## 3. Clone the repository

Go expects it's package to conform to its directory structure.  It reflects the
repositories the packages come from.  So we first create the parent directory:

```
mkdir -p $GOPATH/src/github.com/FabLabBerlin
cd $GOPATH/src/github.com/FabLabBerlin
```

And then clone the Localmachines repository:

```
git clone git@github.com:FabLabBerlin/localmachines.git
```

## 4. Install the Go dependencies

```
cd localmachines
go get
```

## 5. Initialize database

There are many ways to install MySQL on OS X.  The way covered here requires
you to start the MySQL server manually when you need it.  First download
the DMG file [here](https://dev.mysql.com/downloads/mysql/).  Install it as
usual.

Afterwards you can start MySQL with:

```
sudo /usr/local/mysql/support-files/mysql.server start
```

Stop it with:

```
sudo /usr/local/mysql/support-files/mysql.server stop
```

Restart it with:

```
sudo /usr/local/mysql/support-files/mysql.server restart
```

After starting it, you can create a user and an empty database:

```
mysql -u root -p
CREATE USER 'user'@'localhost' IDENTIFIED BY 'pass';
CREATE DATABASE fabsmith;
GRANT ALL PRIVILEGES ON fabsmith.* TO 'user'@'%' WITH GRANT OPTION;
```

Now fill it with data:

```
mysql -u user -p fabsmith < fabsmith.sql
bee migrate -conn="user:pass@tcp(127.0.0.1:3306)/fabsmith"
```

## 6. Setup the backend server

```
cp conf/app.example.conf conf/app.conf
```

Now open conf/app.conf in an editor and setup the parameters as needed.
Usually you want to set

```
runmode=dev
```

unless you are going to prepare a production deployment.

## 7. Install Node.js

You can download a .pkg file from the
[Node.js Download Page](https://nodejs.org/en/download/).

## 8. Install bower and grunt

```
npm -g install bower
npm -g install grunt
```

## 9. Setup React frontend

```
cd clients/machines
npm install
cd $GOPATH/src/github.com/FabLabBerlin/localmachines
```

## 9a. Setup Admin Panel

```
cd clients/admin
npm install
bower install
cd $GOPATH/src/github.com/FabLabBerlin/localmachines
```

## 9b. Setup Signup Page

```
cd clients/signup
npm install
bower install
cd $GOPATH/src/github.com/FabLabBerlin/localmachines
```

## 10. Now you should be able to start the server with:

```
bee run
```

In separate terminal tabs, you should also run:

Terminal Tab 1:

```
cd clients/admin
grunt dev
```

Terminal Tab 2:

```
cd clients/machines
npm run-script dev
```

Terminal Tab 3:

```
cd clients/signup
grunt dev
```

## 10a. Prepare for committing your changes

```
cd $GOPATH/src/github.com/FabLabBerlin/localmachines
./testall
```

The unit tests should run smoothly with no errors.


## 10b. Prepare for deployment

```
cd $GOPATH/src/github.com/FabLabBerlin/localmachines/clients/admin
grunt prod
cd $GOPATH/src/github.com/FabLabBerlin/localmachines/clients/machines
npm run-script prod
cd $GOPATH/src/github.com/FabLabBerlin/localmachines/clients/signup
grunt prod
```

It is okay if you see warnings. Many of them are generated through not so well
maintained 3rd party libraries, that we unfortunately need. :)

Commit the changes. The usual commit message is "app.min.js".