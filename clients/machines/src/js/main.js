var App = require('./components/App');
var FeedbackPage = require('./components/Feedback/FeedbackPage');
var FeedbackStore = require('./stores/FeedbackStore');
var getters = require('./getters');
var LoginChooser = require('./components/Login/LoginChooser');
var LoginStore = require('./stores/LoginStore');
var MachinePage = require('./components/MachinePage/MachinePage');
var MachineStore = require('./stores/MachineStore');
var React = require('react');
var reactor = require('./reactor');
var Router = require('react-router');
var ScrollNavStore = require('./stores/ScrollNavStore');
var SpendingsPage = require('./components/UserProfile/SpendingsPage');
var UserPage = require('./components/UserProfile/UserPage');
var UserStore = require('./stores/UserStore');
var {DefaultRoute, Route, Routes, NotFoundRoute} = require('react-router');
var CategoriesPage = require('./components/Categories/CategoriesPage')
var CategoryPage = require('./components/Categories/CategoryPage')


/*
 * Style dependencies for webpack
 */
require('bootstrap-less');
require('../assets/less/main.less');
require('../assets/less/common.less');
require('font-awesome-webpack');
require('toastr/build/toastr.min.css');
require('vex/css/vex.css');

// Use this to simulate NFC browswer
var debugNfc = false;
if (debugNfc) {
  window.libnfc = {
    debug: true,
    cardRead: {
      connect: function() {},
      disconnect: function() {}
    },
    cardReaderError: {
      connect: function() {},
      disconnect: function() {}
    },
    asyncScan: function() {}
  };
}

/*
 * Defined all the routes of the panel
 */
let routes = (
  <Route name="app" path="/" handler={App} >
    <Route name="login" handler={LoginChooser} />
    <Route name="machine" handler={MachinePage} />
    <Route name="profile" handler={UserPage} />
    <Route name="spendings" handler={SpendingsPage} />
    <Route name="feedback" handler={FeedbackPage} />
    <Route name="categories" handler={CategoriesPage} />
    <Route name="categories/*" handler={CategoryPage} />
    <DefaultRoute handler={CategoriesPage} />
  </Route>
);

/*
 * Define the stores
 */

reactor.registerStores({
  feedbackStore: FeedbackStore,
  loginStore: LoginStore,
  machineStore: MachineStore,
  scrollNavStore: ScrollNavStore,
  userStore: UserStore
});

/*
 * Render everything in the the body of index.html
 */
Router.run(routes, Router.HashLocation, function(Handler) {
  React.render(<Handler />, document.body);
});
