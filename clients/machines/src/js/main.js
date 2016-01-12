var App = require('./components/App');
var FeedbackPage = require('./components/Feedback/FeedbackPage');
var FeedbackStore = require('./modules/Feedback/stores/FeedbackStore');
var getters = require('./getters');
var GlobalStore = require('./modules/Global/stores/GlobalStore');
var LoginChooser = require('./components/Login/LoginChooser');
var LoginStore = require('./modules/Login/stores/LoginStore');
var MachinePage = require('./components/MachinePage/MachinePage');
var MachineStore = require('./modules/Machine/stores/MachineStore');
var React = require('react');
var reactor = require('./reactor');
var ReservationsPage = require('./components/Reservations/ReservationsPage');
var ReservationsStore = require('./modules/Reservations/stores/ReservationsStore');
var ReservationRulesStore = require('./modules/Reservations/stores/ReservationRulesStore');
var Router = require('react-router');
var ScrollNavStore = require('./modules/ScrollNav/stores/ScrollNavStore');
var SpendingsPage = require('./components/UserProfile/SpendingsPage');
var TutoringsStore = require('./modules/Tutorings/stores/TutoringsStore');
var UserPage = require('./components/UserProfile/UserPage');
var UserStore = require('./modules/User/stores/UserStore');
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

// Use this to simulate NFC browswer
var debugNfc = false;
if (debugNfc) {
  window.libnfc = {
    debug: true,
    cardRead: {
      connect() {},
      disconnect() {}
    },
    cardReaderError: {
      connect() {},
      disconnect() {}
    },
    asyncScan() {}
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
    <Route name="reservations" handler={ReservationsPage} />
    <Route name="feedback" handler={FeedbackPage} />
    <DefaultRoute handler={MachinePage} />
  </Route>
);

/*
 * Define the stores
 */

reactor.registerStores({
  feedbackStore: FeedbackStore,
  globalStore: GlobalStore,
  loginStore: LoginStore,
  machineStore: MachineStore,
  reservationsStore: ReservationsStore,
  reservationRulesStore: ReservationRulesStore,
  scrollNavStore: ScrollNavStore,
  tutoringsStore: TutoringsStore,
  userStore: UserStore
});

/*
 * Render everything in the the body of index.html
 */
Router.run(routes, Router.HashLocation, function(Handler) {
  React.render(<Handler />, document.body);
});
