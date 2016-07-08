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
var MachinePage = require('./components/MachinePage/MachinePage');
var Machines = require('./modules/Machines');
var React = require('react');
var reactor = require('./reactor');
var RegisterExisting = require('./components/RegisterExisting');
var ReservationsPage = require('./components/Reservations/ReservationsPage');
var ReservationsStore = require('./stores/ReservationsStore');
var ReservationRulesStore = require('./stores/ReservationRulesStore');
var ScrollNavStore = require('./stores/ScrollNavStore');
var Page = require('./components/UserProfile/SpendingsPage');
var SettingsStore = require('./modules/Settings/stores/store');
var SpendingsPage = require('./components/UserProfile/SpendingsPage');
var TutoringsStore = require('./stores/TutoringsStore');
var UserPage = require('./components/UserProfile/UserPage');
var Users = require('./modules/Users');
var UserStore = require('./stores/UserStore');
var LocationStore = require('./stores/LocationStore');

import { render } from 'react-dom';
import {DefaultRoute, Route, Router, browserHistory} from 'react-router';

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
  tutoringsStore: TutoringsStore,
  userStore: UserStore,
  usersStore: Users.store,
  locationStore: LocationStore
});


const About = React.createClass({
  render() {
    return (
      <div>Aboot</div>
    );
  }
});

var LoaderLocal = require('./components/LoaderLocal');
/*
 * Render everything in the the body of index.html
 */

render((
  <Router history={browserHistory}>
    <Route path="/" component={App} >
      <Route path="machines" component={MachinePage}/>
    </Route>
  </Router>
), document.getElementById('app-container'));

/*render((
  <Router history={browserHistory}>
    <Route path="/" component={App} >
      <Route name="admin" path="admin">
        <Route path="machines" component={AdminMachines} />
        <Route path="machines/:machineId" component={AdminMachine} />
        <Route path="invoices" component={AdminInvoices} />
        <Route path="settings" component={AdminSettings} />
        <Route path="users" component={AdminUsers} />
        <Route path="users/:userId" component={AdminUser} />
      </Route>
      <Route name="forgot_password" path="forgot_password">
        <Route name="email_sent" component={ForgotPassword.EmailSent} />
        <Route name="start" component={ForgotPassword.Start} />
        <Route name="recover" component={ForgotPassword.Recover} />
        <Route name="reset" component={ForgotPassword.Reset} />
        <Route name="done" component={ForgotPassword.Done} />
      </Route>
      <Route name="login" component={LoginChooser} />
      <Route name="machine" component={MachinePage} />
      <Route name="profile" component={UserPage} />
      <Route name="register_existing" component={RegisterExisting} />
      <Route name="spendings" component={SpendingsPage} />
      <Route name="reservations" component={ReservationsPage} />
      <Route name="feedback" component={FeedbackPage} />
      <DefaultRoute component={MachinePage} />
    </Route>
  </Router>
), document.getElementById('app-container'));*/
