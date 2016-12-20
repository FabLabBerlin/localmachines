import { render } from 'react-dom';
import {DefaultRoute, Route, Router, hashHistory, NoRoute} from 'react-router';

if (window.location.pathname === '/product' || window.location.pathname === '/product/') {
  window.location.href = '/machines/#/product';
} else {

var AdminInvoice = require('./components/Admin/Invoices/Invoice');
var AdminInvoices = require('./components/Admin/Invoices/Invoices');
var AdminLocations = require('./components/Admin/Locations/Locations');
var AdminMachine = require('./components/Admin/Machines/Machine');
var AdminMachines = require('./components/Admin/Machines/Machines');
var AdminMembership = require('./components/Admin/Memberships/Membership');
var AdminMemberships = require('./components/Admin/Memberships');
var AdminSettings = require('./components/Admin/Settings/Settings');
var AdminUser = require('./components/Admin/Users/User');
var AdminUsers = require('./components/Admin/Users/Users');
var App = require('./components/App');
var CategoriesStore = require('./modules/Categories/stores/store');
var FeedbackPage = require('./components/Feedback/FeedbackPage');
var FeedbackStore = require('./stores/FeedbackStore');
var ForgotPassword = require('./components/ForgotPassword');
var ForgotPasswordStore = require('./modules/ForgotPassword/stores/store');
var getters = require('./getters');
var GlobalStore = require('./stores/GlobalStore');
var Invoices = require('./modules/Invoices');
var Login = require('./components/Login/Login');
var LoginStore = require('./stores/LoginStore');
var MachinePage = require('./components/Machines/MachinePage');
var MachineInfosPage = require('./components/Machines/MachinePage/Infos');
var MachineReservationPage = require('./components/Machines/MachinePage/Reservations');
var Machines = require('./modules/Machines');
var MachinesPage = require('./components/Machines/Machines');
var Memberships = require('./modules/Memberships');
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
var Location = require('./modules/Location');

/*
 * Style dependencies for webpack
 */
require('bootstrap-less');
require('../assets/less/main.less');
require('../assets/less/common.less');
require('font-awesome-webpack');
require('react-datetime');
require('toastr/build/toastr.min.css');
require('vex/css/vex.css');

/*
 * Define the stores
 */

reactor.registerStores({
  categoriesStore: CategoriesStore,
  feedbackStore: FeedbackStore,
  forgotPasswordStore: ForgotPasswordStore,
  globalStore: GlobalStore,
  invoicesStore: Invoices.store,
  loginStore: LoginStore,
  machineStore: Machines.store,
  membershipsStore: Memberships.store,
  reservationsStore: ReservationsStore,
  reservationRulesStore: ReservationRulesStore,
  scrollNavStore: ScrollNavStore,
  settingsStore: SettingsStore,
  userStore: UserStore,
  usersStore: Users.store,
  locationStore: Location.store,
  locationEditStore: Location.editStore
});


var LoaderLocal = require('./components/LoaderLocal');
/*
 * Render everything in the the body of index.html
 */

render((
  <Router history={hashHistory}>
    <Route path="/" component={App} >
      <Route path="admin">
        <Route path="invoices" component={AdminInvoices} />
        <Route path="invoices/:invoiceId" component={AdminInvoice} />
        <Route path="locations" component={AdminLocations} />
        <Route path="machines" component={AdminMachines} />
        <Route path="machines/:machineId" component={AdminMachine} />
        <Route path="memberships" component={AdminMemberships} />
        <Route path="memberships/:membershipId" component={AdminMembership} />
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
      <Route path="login" component={Login} />
      <Route path="machines" component={MachinesPage} />
      <Route path="machines/:machineId" component={MachinePage} />
      <Route path="machines/:machineId/infos" component={MachineInfosPage} />
      <Route path="machines/:machineId/reservations" component={MachineReservationPage} />
      <Route path="product" component={ProductPage} />
      <Route path="profile" component={UserPage} />
      <Route path="register_existing" component={RegisterExisting} />
      <Route path="spendings" component={SpendingsPage} />
      <Route path="reservations" component={ReservationsPage} />
      <Route path="feedback" component={FeedbackPage} />
      <Route path="/" component={MachinesPage} />
    </Route>
  </Router>
), document.getElementById('app-container'));

} // if ( ... /product ... )
