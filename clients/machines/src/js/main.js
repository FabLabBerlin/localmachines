import { render } from 'react-dom';
import {DefaultRoute, Route, Router, hashHistory, NoRoute} from 'react-router';

if (window.location.pathname === '/product' || window.location.pathname === '/product/') {
  window.location.href = '/machines/#/product';
} else {

var AdminInvoices = require('./components/Admin/Invoices/Invoices');
var AdminMachine = require('./components/Admin/Machines/Machine');
var AdminMachines = require('./components/Admin/Machines/Machines');
var AdminSettings = require('./components/Admin/Settings/Settings');
var AdminUser = require('./components/Admin/Users/User');
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
var MachineNewPage = require('./components/MachinesNew/MachinePage');
var MachineNewReservationPage = require('./components/MachinesNew/MachinePage/Reservations');
var MachinePage = require('./components/MachinePage/MachinePage');
var Machines = require('./modules/Machines');
var MachinesNewPage = require('./components/MachinesNew/Machines');
var React = require('react');
var ProductPage = require('./components/ProductPage/ProductPage');
var reactor = require('./reactor');
var RegisterExisting = require('./components/RegisterExisting');
var ReservationsPage = require('./components/Reservations/ReservationsPage');
var ReservationsStore = require('./stores/ReservationsStore');
var ReservationRulesStore = require('./stores/ReservationRulesStore');
var ScrollNavStore = require('./stores/ScrollNavStore');
var Page = require('./components/UserProfile/SpendingsPage');
var SettingsStore = require('./modules/Settings/stores/store');
var SpendingsPage = require('./components/UserProfile/SpendingsPage');
var UserPage = require('./components/UserProfile/UserPage');
var Users = require('./modules/Users');
var UserStore = require('./stores/UserStore');
var LocationStore = require('./stores/LocationStore');

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
  userStore: UserStore,
  usersStore: Users.store,
  locationStore: LocationStore
});


var LoaderLocal = require('./components/LoaderLocal');
/*
 * Render everything in the the body of index.html
 */

render((
  <Router history={hashHistory}>
    <Route path="/" component={App} >
      <Route path="admin">
        <Route path="machines" component={AdminMachines} />
        <Route path="machines/:machineId" component={AdminMachine} />
        <Route path="invoices" component={AdminInvoices} />
        <Route path="settings" component={AdminSettings} />
        <Route path="users" component={AdminUsers} />
        <Route path="users/:userId" component={AdminUser} />
      </Route>
      <Route path="forgot_password">
        <Route path="email_sent" component={ForgotPassword.EmailSent} />
        <Route path="start" component={ForgotPassword.Start} />
        <Route path="recover" component={ForgotPassword.Recover} />
        <Route path="reset" component={ForgotPassword.Reset} />
        <Route path="done" component={ForgotPassword.Done} />
      </Route>
      <Route path="login" component={LoginChooser} />
      <Route path="machine" component={MachinePage} />
      <Route path="machines" component={MachinesNewPage} />
      <Route path="machines/:machineId" component={MachineNewPage} />
      <Route path="machines/:machineId/reservations" component={MachineNewReservationPage} />
      <Route path="product" component={ProductPage} />
      <Route path="profile" component={UserPage} />
      <Route path="register_existing" component={RegisterExisting} />
      <Route path="spendings" component={SpendingsPage} />
      <Route path="reservations" component={ReservationsPage} />
      <Route path="feedback" component={FeedbackPage} />
      <Route path="/" component={MachinePage} />
    </Route>
  </Router>
), document.getElementById('app-container'));

} // if ( ... /product ... )
