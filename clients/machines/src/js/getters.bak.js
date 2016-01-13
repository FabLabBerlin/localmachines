var _ = require('lodash');
var { Day, Time } = require('./components/Reservations/helpers');
var { toCents, subtractVAT } = require('./components/UserProfile/helpers');
var moment = require('moment');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


/*
 * Global state related getters
 */
const getIsLoading = [
  ['globalStore'],
  (machineStore) => {
    return machineStore.get('loading');
  }
];




/*
 * Spendings (Page) related getters
 */







export default {
  getIsLogged, getUid, getFirstTry, getLoginSuccess, getLastActivity,
  getUser, getActivations, getMachines, getMachinesById, getMachineUsers, getIsLoading, getBill, getBillMonths, getMonthlyBills, getMemberships, getMembershipsByMonth,
  getFeedbackSubject, getFeedbackSubjectDropdown, getFeedbackSubjectOtherText, getFeedbackMessage,
  getScrollUpEnabled, getScrollDownEnabled, getScrollPosition
};
