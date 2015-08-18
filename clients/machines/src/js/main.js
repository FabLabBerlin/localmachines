var App = require('./components/App');
var getters = require('./getters');
var LoginChooser = require('./components/LoginChooser');
var LoginStore = require('./stores/LoginStore');
var MachinePage = require('./components/MachinePage');
var MachineStore = require('./stores/MachineStore');
var React = require('react');
var reactor = require('./reactor');
var Router = require('react-router');
var {DefaultRoute, Route, Routes, NotFoundRoute} = require('react-router');

/*
 * Style dependencies for webpack
 */
require('bootstrap-less');
require('../assets/less/main.less');
require('../assets/less/common.less');
require('font-awesome-webpack');
require('toastr/build/toastr.min.css');
require('vex/css/vex.css');


/*
 * Defined all the routes of the panel
 */
let routes = (
  <Route name="app" path="/" handler={App} >
    <Route name="machine" handler={MachinePage} />
    <Route name="login" handler={LoginChooser} />
    <DefaultRoute handler={MachinePage} />
  </Route>
);

/*
 * Define the stores
 */

reactor.registerStores({
  loginStore: LoginStore,
  machineStore: MachineStore
});

/*
 * Render everything in the the body of index.html
 */
Router.run(routes, Router.HashLocation, function(Handler) {
  React.render(<Handler />, document.body);
});
