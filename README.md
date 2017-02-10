# Local Machines

Resource management software for Fab Labs. REST API built with the [BeeGo](http://beego.me) framework. Front end clients currently build with [AngularJS](https://angularjs.org) and [ReactJS](https://facebook.github.io/react/). 

## Table of contents
- [Quick-start](#quick-start)
- [Versioning](#versioning)
- [Runmode](#runmode)
- [Development](#development)
  - [Code Conventions](#code-conventions)
  - [Core Development](#core-development)
  - [Clients Development](#clients-development)
- [API Documentation](#api-documentation)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)
- [TODO](#TODO)

## Quick-Start

Quickly install Localmachines to get it running on Linux or OS X. For more
 specialized/detailed setups, see [here](#specialized-setups).

1. Install [Go](https://golang.org). You can download binaries from the
  [Go Download Page](https://golang.org/dl/). For OS X you can get the
  .pkg file. On Ubuntu you can just enter `sudo apt-get install golang-go`.
2. Set Go environment variables:
   `export GOROOT=/usr/local/go` and
   `export GOPATH=$HOME/go` (you may need `mkdir $HOME/go` beforehands)

   Most users add both lines to their `$HOME/.bash_profile` file.
3. Clone the repository:

   ```
   mkdir -p $GOPATH/src/github.com/FabLabBerlin
   cd $GOPATH/src/github.com/FabLabBerlin
   git clone git@github.com:FabLabBerlin/localmachines.git
   ```
4a. Install Glide (optional)
   ```
   curl https://glide.sh/get | sh
   ```

   The [Glide Github page](https://github.com/Masterminds/glide#install) shows alternative installation methods.

4b. Install the Go dependencies *either* using Glide:
   ```
   cd localmachines
   glide install
   ```

   But you can also install them manually if you wish.

   ImageMagick is also needed.

5. Initialize database
   ```
   mysql -u user -p fabsmith < fabsmith.sql
   bee migrate -conn="root:@tcp(127.0.0.1:3306)/fabsmith"
   ```

6. Setup the backend server:
   ```
   cp conf/app.example.conf conf/app.conf
   ```
7. Install [Node.js LTS](https://nodejs.org/)
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

10. Setup Angular frontends:

   ```
   cd clients/admin
   npm install
   bower install

   cd clients/signup
   npm install
   bower install

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

## Specialized-Setups

### OS X Development Setup

Detailed instruction on how to install Localmachines on OS X are
[here](docs/osx).

### Raspberry Pi

Instructions on how to install Localmachines on a Raspberry Pi are
[here](docs/raspi).

## General structure

### Versioning

localmachines uses versioning scheme that is loosely following the rules of
[Semantic Versioning](http://semver.org): [MAJOR].[MINOR].[PATCH]

### Overall File structure

```
clients/      HTML5 Clients for Admins and End-Users
conf/         Config files
controllers/  Backend Controllers
database/     Database Migrations
docs/         Documentation
files/        Invoices
gateway/      Lab "IoT" Gateway Server
icons/        Icons
lib/          Reusable general purpose functions
models/       Backend Models
netswitch/    NetSwitch TestKit
pinger/       Ping tool
routers/      Backend Routes
scripts/      Scripts
swagger/      Swagger Documentation Generator
tests/        Backend Unit Tests
```

### Admin Client

#### File Structure

```
dev/                         Development files
dev/assets/                  Static Assets
dev/assets/css/              Generated CSS (do not edit)
dev/assets/img/              Images
dev/assets/js/               Misc. Javascript (deprecated)
dev/assets/less/             Less files (edit CSS from here)
dev/ng-components/           External Angular Modules
dev/ng-components/version/
dev/ng-modules/               Angular Modules
dev/ng-modules/activations/   Activations Page
dev/ng-modules/api/           HTTP API helper functions
dev/ng-modules/bookings/      Bookings (obsolete)
dev/ng-modules/coworking/     Co-Working
dev/ng-modules/dashboard/     Dashboard with Metrics
dev/ng-modules/invoices/      Invoices
dev/ng-modules/login/         Login
dev/ng-modules/machine/       Edit a Machine
dev/ng-modules/machines/      Machine List
dev/ng-modules/mainmenu/      Menu in the Header
dev/ng-modules/membership/    Edit a Memberhip
dev/ng-modules/memberships/   Memberships List
dev/ng-modules/priceunit/     Select box for Price Unit
dev/ng-modules/productlist/   Generic Product List Table
dev/ng-modules/randomtoken/   "Are you sure?" generator
dev/ng-modules/reservation/   Edit a Reservation
dev/ng-modules/reservations/  Reservations List
dev/ng-modules/settings/      Settings View
dev/ng-modules/space/         Edit a Space
dev/ng-modules/spacepurchase/ Create a Space Purchase
dev/ng-modules/spaces/        Space List
dev/ng-modules/tutoring/      Tutoring
dev/ng-modules/user/          Create/Edit a User
dev/ng-modules/users/         Users List
prod/                         Production files ("`grunt prod`"")
...
```

#### Writing a new Page

- Create directory `dev/ng-modules/foobar`
- Create a file `dev/ng-modules/foobar/foobar.js`

A basic controller looks like this:

```
(function(){

'use strict';

var app = angular.module('fabsmith.admin.foobar', 
 ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/foobar', {
    templateUrl: 'ng-modules/foobar/foobar.html',
    controller: 'FoobarCtrl'
  });
}]); // app.config

app.controller('FoobarCtrl', ['$scope', '$http', '$location', 'randomToken', 
 function($scope, $http, $location, randomToken) {

$http({
  method: 'GET',
  url: '/api/my_data'
})
.success(function(myData) {
  $scope.myData = myData;
  toastr.info('My Data successfully loaded.');
})
.error(function() {
  toastr.error('Error loading data');
});

$scope.hello = function() {
  toastr.info('Hello World!');
};

}]); // app.controller

})(); // closure
```

- Create a file `dev/ng-modules/foobar/foobar.html`:

```
<div id="admin-global-config" class="container-fluid">

<h1>Foobar</h1>

<h2>My Data: {{myData.Bar}}</h2>

<button class="btn btn-primary btn-lg"
        ng-click="hello()">
  Say hello
</button>

</div> <!-- /container -->
```

- Load the JS file in `dev/index.html`:

Add it after the other module files:

```
...
<script src="/admin/ng-modules/tutoring/tutoring.js"></script>
<script src="/admin/ng-modules/tutoring/tutor.js"></script>
<script src="/admin/ng-modules/tutoring/purchase.js"></script>
<script src="/admin/ng-modules/foobar/foobar.js"></script>
```

- Load the module in `dev/ng-main.js`:

Add it after the other modules:

```
var app = angular.module('fabsmith', [
  'ngRoute',
  'fabsmith.admin.login',
  'fabsmith.admin.api',
  'fabsmith.admin.coworking',
  ...
  'fabsmith.admin.foobar'
]);
```

### End-User "Machines" Client

#### File Structure

The source code is within the `src/` subdirectory:

```
assets/                              Static Assets
assets/less/                         Less files (edit CSS from here)
js/                                  JS code
js/actions/                          Flux Actions
js/actions/__test__/
js/components/                       React Components
js/components/Feedback/              Feedback Form
js/components/Feedback/__test__/
js/components/Login/                 Login Pages
js/components/Login/__test__/
js/components/MachinePage/           Machine Page
js/components/MachinePage/__test__/
js/components/MachinePage/Machine/   Machine in the List
js/components/Reservations/          Reservations Module
js/components/Reservations/__test__/
js/components/UserProfile/           User Profile Page
js/components/UserProfile/__test__/
js/stores/                           Flux Stores
js/stores/__test__/
```

#### Writing a new Module

- Create a new file `js/actions/FooActions.js`:

```
var $ = require('jquery');
var actionTypes = require('../actionTypes');
var reactor = require('../reactor');

var FooActions = {

  loadFoos() {
    $.ajax({
      url: '/api/foos',
      success(foos) {
        reactor.dispatch(actionTypes.SET_FOOS, { foos });
      },
      error(xhr, status, err) {
        toastr.error('Error loading locations');
      }
    });
  }

};

export default FooActions;
```

Note that we make a lot of use of ECMAScript 6 here.  When calling a function
`function foobar({data}) {...}`, the function must be called with
`foobar({data})`.  If you call it with `foobar({baz})`, `data` will be
undefined.  So parameters names must match!

- Declare the new Action Type `SET_FOOS` in `js/actionTypes.js`:

```
var keyMirror = require('react/lib/keyMirror');

export default keyMirror({
  API_GET_LOGOUT: null,
  SUCCESS_LOGIN: null,
...
  SET_FOOS: null
});
```

- Create a Flux store in `js/stores/FooStore.js`:

```
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  foos: undefined
});

var FooStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_FOOS, setFoos);
  }
});

function setFoos(state, { foos }) {
  return state.set('foos', foos);
}

export default FooStore;
```

- Register the Flux store in `js/main.js`:

Add

```
var FooStore = require('./stores/FooStore');
```

to the imports and register the actual store:

```
reactor.registerStores({
  feedbackStore: FeedbackStore,
  globalStore: GlobalStore,
  loginStore: LoginStore,
  machineStore: MachineStore,
  reservationsStore: ReservationsStore,
  reservationRulesStore: ReservationRulesStore,
  scrollNavStore: ScrollNavStore,
  tutoringsStore: TutoringsStore,
  userStore: UserStore,
  locationStore: LocationStore,
  fooStore: FooStore
});
```

- Write a getter in `js/getters.js`:

```
const getFoos = [
  ['fooStore'],
  (fooStore) => {
    return fooStore.get('foos');
  }
];
```

and export it at the bottom of the file:

```
export default {
  ...
  getTutorings,
  getFoos
};
```

- Write the React Class in `js/components/Foo.js`:

```
var _ = require('lodash');
var FooActions = require('../../actions/FooActions');
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


var FooWidget = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      foos: getters.getFoos
    };
  },

  componentDidMount() {
    FooActions.loadFoos();
  },

  click() {
    toastr.info('Click!');
  },

  render() {
    return (
      <div className="container">
        {
          _.isUndefined(this.state.foos) ?
          (
            <div>Please be patient, loading data...</div>
          ) : (
            <div>
              <h2>Foos</h2>
              {this.state.foos.map((foo, i) => {
                <h3>Foo No. {i} is named {foo.Name}</h3>
              })}
              <button onClick={this.click}>Click me</button>
            </div>
          )
        }
      </div>
    );
  }
});

export default FooWidget;
```

- Now the component can be used in other components:

```
var Foo = require('./Foo');

...
<Foo/>
```

### Backend

#### File Structure

Go uses packages to encapsulate code and data.  In a package, only the
functions and variables with uppercase names are "public", the others are
"private".  Declare functions as private to keep the APIs lean.

You can get an overview with the `godoc` tool:

```
godoc -http=:6060
```

Now navigate with your Web Browser to the
[localmachines pkg on localhost:6060](http://localhost:6060/pkg/github.com/FabLabBerlin/localmachines/).
There you get an interactive overview of the whole Backend code.

## Development

The development environment is a constant work in progress and does not implement a decent test-driven development workflow yet. Soon.

When ready to move the system to `prod` runmode, run `grunt prod` to compile the production mode of the Angular JS applications.

More info about the development workflow will be added to Wiki.

### Code-Conventions

Go files should have a maximum of 80 characters per lines, JS files 100
characters per line maximum.

Local variables should be as short as possible. Public variables and functions
must have terse but quickly comprehendable names.

### Testing

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

1. Compile the documentation by using the following command

  ```
  bee generate docs
  ```

2. Run the server

  ```
  bee run
  ```

3. Access API documentation via: 

  ```
  http://localhost:8080/swagger/
  ```

## Troubleshooting

### Too many open files

Beego watches a lot of files to automatically re-compile the code. On OS X the maximum number of open files is quite low by default. See here how to solve this: http://unix.stackexchange.com/a/221988

## TODO

- Merge functionality of clients/admin to clients/machines
- Merge functionality of clients/signup to clients/machines
- Flatten directory structure in clients
