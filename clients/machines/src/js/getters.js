var _ = require('lodash');
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
const getUserInfo = [
  ['userStore'],
  (userStore) => {
    return userStore.get('userInfo');
  }
];

/*
 * Spendings (Page) related getters
 */
const getBillInfo = [
  ['userStore'],
  (userStore) => {
    return userStore.get('billInfo');
  }
];

const getBillMonths = [
  getUserInfo,
  (userInfo) => {
    var months = [];
    var created = moment(userInfo.get('Created'));
    if (created.unix() <= 0) {
      created = moment('2015-07-01');
    }
    for (var t = created; t.isBefore(moment()); t.add(1, 'd')) {
      months.push(t.clone());
    }
    months = _.uniq(months, function(month) {
      return month.format('MMM YYYY');
    });
    return months;
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
      var end = moment(membership.EndDate);
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

const getMonthlyBills = [
  getBillInfo,
  getBillMonths,
  getMembershipsByMonth,
  (billInfo, billMonths, membershipsByMonth) => {
    if (!billInfo) {
      return undefined;
    }
    var activations = billInfo.Activations;
    var activationsByMonth = _.groupBy(activations, function(info) {
      return moment(info.TimeStart).format('MMM YYYY');
    });
    var monthlyBills = _.map(billMonths, function(m) {
      var month = m.format('MMM YYYY');
      var monthlyBill = {
        month: month,
        activations: [],
        memberships: [],
        sums: {
          activations: {
            priceInclVAT: 0,
            priceExclVAT: 0,
            priceVAT: 0,
            durations: 0
          },
          memberships: {
            priceInclVAT: 0,
            priceExclVAT: 0,
            priceVAT: 0
          },
          total: {}
        }
      };

      /*
       * Collect activations and sum for the totals
       */
      _.eachRight(activationsByMonth[month], function(info) {
        
        var duration = moment.duration(moment(info.TimeEnd)
          .diff(moment(info.TimeStart))).asSeconds();
        
        monthlyBill.sums.durations += duration;
        var priceInclVAT = toCents(info.DiscountedTotal);
        var priceExclVAT = toCents(subtractVAT(info.DiscountedTotal));
        var priceVAT = priceInclVAT - priceExclVAT;
        monthlyBill.sums.activations.priceInclVAT += priceInclVAT;
        monthlyBill.sums.activations.priceExclVAT += priceExclVAT;
        monthlyBill.sums.activations.priceVAT += priceVAT;
        
        monthlyBill.activations.push({
          MachineName: info.Machine.Name,
          TimeStart: moment(info.TimeStart),
          duration: duration,
          priceExclVAT: priceExclVAT,
          priceVAT: priceVAT,
          priceInclVAT: priceInclVAT
        });
      });

      /*
       * Collect memberships for each month
       */
      _.each(membershipsByMonth[month], function(membership) {
        var totalPrice = toCents(membership.MonthlyPrice);
        var priceExclVat = toCents(subtractVAT(membership.MonthlyPrice));
        var vat = totalPrice - priceExclVat;
        monthlyBill.memberships.push({
          startDate: moment(membership.StartDate),
          endDate: moment(membership.EndDate),
          priceExclVAT: priceExclVat,
          priceVAT: vat,
          priceInclVAT: totalPrice
        });
        monthlyBill.sums.memberships.priceInclVAT += totalPrice;
        monthlyBill.sums.memberships.priceExclVAT += priceExclVat;
        monthlyBill.sums.memberships.priceVAT += vat;
      });

      monthlyBill.sums.total = {
        priceInclVAT: monthlyBill.sums.activations.priceInclVAT + monthlyBill.sums.memberships.priceInclVAT,
        priceExclVAT: monthlyBill.sums.activations.priceExclVAT + monthlyBill.sums.memberships.priceExclVAT,
        priceVAT: monthlyBill.sums.activations.priceVAT + monthlyBill.sums.memberships.priceVAT
      };

      return monthlyBill;
    });
    monthlyBills.reverse();
    return monthlyBills;
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

const getMachinesById = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('machinesById') || toImmutable({});
  }
];

const getMachineInfo = [
  getMachinesById,
  (machinesById) => {
    return machinesById.sortBy(m => m.Name);
  }
];

const getMachineUsers = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('machineUsers') || {};
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
 * Reservations related getters
 */
const getNewReservation = [
  ['reservationsStore'],
  (reservationsStore) => {
    return reservationsStore.get('create');
  }
];

const getNewReservationTimes = [
  ['reservationsStore'],
  (reservationsStore) => {
    return reservationsStore.get('create').get('times');
  }
];

const getNewReservationFrom = [
  getNewReservationTimes,
  (reservationTimes) => {
    console.log('reservationTimes.toArray():', reservationTimes.toArray());
    var selectedTimes = _.filter(reservationTimes.toArray(), function(t) {
      return t.get('selected');
    });
    console.log('selectedTimes:', selectedTimes);
    if (selectedTimes.length > 0) {
      return _.first(selectedTimes).get('start').toDate();
    }
  }
];

const getNewReservationTo = [
  getNewReservationTimes,
  (reservationTimes) => {
    var selectedTimes = _.filter(reservationTimes.toArray(), function(t) {
      return t.get('selected');
    });
    if (selectedTimes.length > 0) {
      return _.last(selectedTimes).get('end').toDate();
    }
  }
];

const getReservations = [
  ['reservationsStore'],
  (reservationsStore) => {
    return reservationsStore.get('reservations');
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
  getUserInfo, getActivationInfo, getMachineInfo, getMachinesById, getMachineUsers, getIsLoading, getBillInfo, getBillMonths, getMonthlyBills, getMembership, getMembershipsByMonth,
  getFeedbackSubject, getFeedbackSubjectDropdown, getFeedbackSubjectOtherText, getFeedbackMessage,
  getNewReservation, getNewReservationTimes, getNewReservationFrom, getNewReservationTo, getReservations,
  getScrollUpEnabled, getScrollDownEnabled, getScrollPosition
};
