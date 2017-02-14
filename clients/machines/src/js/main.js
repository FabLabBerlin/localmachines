import { render } from 'react-dom';
import {DefaultRoute, Route, Router, hashHistory, NoRoute} from 'react-router';

import AdminCategories from './components/Admin/Categories/Categories';
import AdminInvoice from './components/Admin/Invoices/Invoice';
import AdminInvoices from './components/Admin/Invoices/Invoices';
import AdminLocations from './components/Admin/Locations/Locations';
import AdminMachine from './components/Admin/Machines/Machine';
import AdminMachines from './components/Admin/Machines/Machines';
import AdminMembership from './components/Admin/Memberships/Membership';
import AdminMemberships from './components/Admin/Memberships';
import AdminSettings from './components/Admin/Settings/Settings';
import AdminUser from './components/Admin/Users/User';
import AdminUsers from './components/Admin/Users/Users';
import App from './components/App';
import CategoriesStore from './modules/Categories/stores/store';
import FeedbackPage from './components/Feedback/FeedbackPage';
import FeedbackStore from './stores/FeedbackStore';
import ForgotPassword from './components/ForgotPassword';
import ForgotPasswordStore from './modules/ForgotPassword/stores/store';
import getters from './getters';
import GlobalStore from './stores/GlobalStore';
import Invoices from './modules/Invoices';
import LoaderLocal from './components/LoaderLocal';
import Login from './components/Login/Login';
import LoginStore from './stores/LoginStore';
import MachinePage from './components/Machines/MachinePage';
import MachineInfosPage from './components/Machines/MachinePage/Infos';
import MachineReservationPage from './components/Machines/MachinePage/Reservations';
import Machines from './modules/Machines';
import MachinesPage from './components/Machines/Machines';
import Memberships from './modules/Memberships';
import React from 'react';
import ProductPage from './components/ProductPage/ProductPage';
import reactor from './reactor';
import RegisterExisting from './components/RegisterExisting';
import ReservationsPage from './components/Reservations/ReservationsPage';
import ReservationsStore from './stores/ReservationsStore';
import ReservationRulesStore from './stores/ReservationRulesStore';
import ScrollNavStore from './stores/ScrollNavStore';
import Page from './components/UserProfile/SpendingsPage';
import SettingsStore from './modules/Settings/stores/store';
import SpendingsPage from './components/UserProfile/SpendingsPage';
import UserPage from './components/UserProfile/UserPage';
import Users from './modules/Users';
import UserStore from './stores/UserStore';
import Location from './modules/Location';

if (window.location.pathname === '/product' || window.location.pathname === '/product/') {
  window.location.href = '/machines/#/product';
} else {

/*
 * Style dependencies for webpack
 */
require('bootstrap-less');
require('../assets/less/main.less');
require('../assets/less/common.less');
require('font-awesome/css/font-awesome.css');
require('react-datetime');
require('toastr/build/toastr.min.css');
require('vex/css/vex.css');

/*
 * Define the stores
 */

_.each({
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
}, (s, k) => {
  console.log('store[' + k + ']=', s)
});

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


/*
 * Render everything in the the body of index.html
 */

render((
  <Router history={hashHistory}>
    <Route path="/" component={App} >
      <Route path="admin">
        <Route path="categories" component={AdminCategories} />
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
