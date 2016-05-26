var AdminInvoices = require('./components/Admin/Invoices/Invoices');
var AdminMachine = require('./components/Admin/Machines/Machine');
var AdminMachines = require('./components/Admin/Machines/Machines');
var AdminSettings = require('./components/Admin/Settings/Settings');
var AdminUsers = require('./components/Admin/Users/Users');
var App = require('./components/App');
var FeedbackPage = require('./components/Feedback/FeedbackPage');
var FeedbackStore = require('./stores/FeedbackStore');
var ForgotPassword = require('./components/ForgotPassword');
var ForgotPasswordStore = require('./modules/ForgotPassword/stores/store');
var getters = require('./getters');
var GlobalStore = require('./stores/GlobalStore');
var Invoices = require('./modules/Invoices');
var LoginChooser = require('./components/Login/LoginChooser');
var LoginStore = require('./stores/LoginStore');
var MachinePage = require('./components/MachinePage/MachinePage');
var Machines = require('./modules/Machines');
var React = require('react');
var reactor = require('./reactor');
var RegisterExisting = require('./components/RegisterExisting');
var ReservationsPage = require('./components/Reservations/ReservationsPage');
var ReservationsStore = require('./stores/ReservationsStore');
var ReservationRulesStore = require('./stores/ReservationRulesStore');
var Router = require('react-router');
var ScrollNavStore = require('./stores/ScrollNavStore');
var Page = require('./components/UserProfile/SpendingsPage');
var SettingsStore = require('./modules/Settings/stores/store');
var SpendingsPage = require('./components/UserProfile/SpendingsPage');
var TutoringsStore = require('./stores/TutoringsStore');
var UserPage = require('./components/UserProfile/UserPage');
var Users = require('./modules/Users');
var UserStore = require('./stores/UserStore');
var LocationStore = require('./stores/LocationStore');

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
    <Route name="admin" path="admin">
      <Route path="machines" handler={AdminMachines} />
      <Route path="machines/:machineId" handler={AdminMachine} />
      <Route path="invoices" handler={AdminInvoices} />
      <Route path="settings" handler={AdminSettings} />
      <Route path="users" handler={AdminUsers} />
    </Route>
    <Route name="forgot_password" path="forgot_password">
      <Route name="email_sent" handler={ForgotPassword.EmailSent} />
      <Route name="start" handler={ForgotPassword.Start} />
      <Route name="recover" handler={ForgotPassword.Recover} />
      <Route name="reset" handler={ForgotPassword.Reset} />
      <Route name="done" handler={ForgotPassword.Done} />
    </Route>
    <Route name="login" handler={LoginChooser} />
    <Route name="machine" handler={MachinePage} />
    <Route name="profile" handler={UserPage} />
    <Route name="register_existing" handler={RegisterExisting} />
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
  forgotPasswordStore: ForgotPasswordStore,
  globalStore: GlobalStore,
  invoicesStore: Invoices.store,
  loginStore: LoginStore,
  machineStore: Machines.store,
  reservationsStore: ReservationsStore,
  reservationRulesStore: ReservationRulesStore,
  scrollNavStore: ScrollNavStore,
  settingsStore: SettingsStore,
  tutoringsStore: TutoringsStore,
  userStore: UserStore,
  usersStore: Users.store,
  locationStore: LocationStore
});

/*
 * Render everything in the the body of index.html
 */
Router.run(routes, Router.HashLocation, function(Handler) {
  React.render(<Handler />, document.body);
});
