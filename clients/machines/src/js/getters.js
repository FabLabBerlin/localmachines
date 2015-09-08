var _ = require('lodash');
var moment = require('moment');
var reactor = require('./reactor');


/*
 * Login state related getters
 */
const getIsLogged = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('isLogged');
  }
];

const getUid = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('uid');
  }
];

const getFirstTry = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('firstTry');
  }
];

const getLoginSuccess = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('loginSuccess');
  }
];

const getLastActivity = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('lastActivity');
  }
];


/*
 * User (Page) related getters
 */
 const getBillInfo = [
  ['userStore'],
  (userStore) => {
    return userStore.get('billInfo');
  }
 ];

 const getMembership = [
  ['userStore'],
  (userStore) => {
    return userStore.get('membershipInfo');
  }
 ];

 const getMembershipsByMonth = [
  ['userStore'],
  (userStore) => {
    var byMonths = {};
    _.each(userStore.get('membershipInfo'), function(membership) {
      var start = moment(membership.StartDate);
      var end = moment(membership.StartDate).add(membership.Duration, 'd');
      for (var t = start; t.isBefore(end); t = t.add(1, 'M')) {
        var month = t.format('MMM YYYY');
        if (!byMonths[month]) {
          byMonths[month] = [];
        }
        byMonths[month].push(membership);
      }
    });
    return byMonths;
  }
 ];


 const getUserInfo = [
  ['userStore'],
  (userStore) => {
    return userStore.get('userInfo');
  }
];


/*
 * Machine (Page) related getters
 */
const getActivationInfo = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('activationInfo');
  }
];

const getMachineInfo = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('machineInfo') || [];
  }
];

const getMachine = function(id) {
  return reactor.evaluate([
    ['machineStore'],
    (machineStore) => {
      var machines = machineStore.get('machineInfo') || [];
      var machine;
      _.each(machines, function(m) {
        if (m.Id === id) {
          machine = m;
          return false;
        }
      });
      return machine;
    }
  ]);
};

const getMachineUsers = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('machineUsers') || {};
  }
];

const getIsLoading = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('loading');
  }
];


/*
 * Feedback related getters
 */
const getFeedbackSubject = [
  ['feedbackStore'],
  (feedbackStore) => {
    if (feedbackStore.get('subject-dropdown') === 'Other') {
      return feedbackStore.get('subject-other-text');
    } else {
      return feedbackStore.get('subject-dropdown');
    }
  }
];

const getFeedbackSubjectDropdown = [
  ['feedbackStore'],
  (feedbackStore) => {
    return feedbackStore.get('subject-dropdown');
  }
];

const getFeedbackSubjectOtherText = [
  ['feedbackStore'],
  (feedbackStore) => {
    return feedbackStore.get('subject-other-text');
  }
];

const getFeedbackMessage = [
  ['feedbackStore'],
  (feedbackStore) => {
    return feedbackStore.get('message');
  }
];


/*
 * Scroll Navigation related getters
 */
const getScrollUpEnabled = [
  ['scrollNavStore'],
  (scrollNavStore) => {
    return scrollNavStore.get('upEnabled');
  }
];

const getScrollDownEnabled = [
  ['scrollNavStore'],
  (scrollNavStore) => {
    return scrollNavStore.get('downEnabled');
  }
];

const getScrollPosition = [
  ['scrollNavStore'],
  (scrollNavStore) => {
    return scrollNavStore.get('position');
  }
];

export default {
  getIsLogged, getUid, getFirstTry, getLoginSuccess, getLastActivity,
  getUserInfo, getActivationInfo, getMachineInfo, getMachine, getMachineUsers, getIsLoading, getBillInfo, getMembership, getMembershipsByMonth,
  getFeedbackSubject, getFeedbackSubjectDropdown, getFeedbackSubjectOtherText, getFeedbackMessage,
  getScrollUpEnabled, getScrollDownEnabled, getScrollPosition
};
