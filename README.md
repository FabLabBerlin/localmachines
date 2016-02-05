# Localmachines
Machine activation software for Fab Labs. Build with [BeeGo](http://beego.me) framework for [GoLang](https://golang.org) and [Angular.js](https://angularjs.org).

## Table of contents
- [Quick-start](#quick-start)
- [Versioning](#versioning)
- [Runmode](#runmode)
- [Development](#development)
  - [Core Development](#core-development)
  - [Clients Development](#clients-development)
- [API Documentation](#api-documentation)

##Quick-Start
Quickly install Localmachines to get it running on Linux or OS X. For more
 specialized/detailed setups, see [here](#specialized-setups).

1. Install [Go](https://golang.org). You can download binaries from the
  [Go Download Page](https://golang.org/dl/). For OS X you can get the
  .pkg file. On Ubuntu you can just enter `sudo apt-get install golang-go`.
2. Set Go environment variables:
   `export GOROOT=/usr/local/go` and
   `export GOPATH=$HOME/go` (you may need `mkdir $HOME/go` beforehands)
3. Clone the repository:

   ```
   mkdir -p $GOPATH/src/github.com/FabLabBerlin
   cd $GOPATH/src/github.com/FabLabBerlin
   git clone git@github.com:FabLabBerlin/localmachines.git
   ```
4. Install the Go dependencies:
   ```
   cd localmachines
   go get
   ```

5. Initialize database
   ```
   mysql -u user -p fabsmith < fabsmith.sql
   bee migrate -conn="root:@tcp(127.0.0.1:3306)/fabsmith"
   ```

6. Setup the backend server:
   ```
   cp conf/app.example.conf conf/app.conf
   ```
7. Install [Node.js](https://nodejs.org/)
8. Install [bower](http://bower.io/) and [grunt](http://gruntjs.com/):

   ```
   npm -g install bower
   npm -g install grunt
   ```

9. Setup React frontend:

   ```
   cd clients/machines
   npm install
   cd $GOPATH/src/github.com/FabLabBerlin/localmachines
   ```

10. Now you should be able to start the server with:

   ```
   bee run
   ```

   or alternatively:

   ```
   go build
   ./localmachines
   ```

##Specialized-Setups

###OS X Development Setup

Detailed instruction on how to install Localmachines on OS X are
[here](docs/osx).

###Raspberry Pi

Instructions on how to install Localmachines on a Raspberry Pi are
[here](docs/raspi).

##Versioning
FabSmith will benefit from semantic versioning. Read about it [here](http://semver.org).

## Development

The development environment is a constant work in progress and does not implement a decent test-driven development workflow yet. Soon.

When ready to move the system to `prod` runmode, run `grunt prod` to compile the production mode of the Angular JS applications.

More info about the development workflow will be added to Wiki.

## Testing

Use `./testall` to run all tests in the project. Make sure that the file
`./tests/conf/app.conf` exists.

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
