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
    var purchases = billInfo.Purchases;
    var purchasesByMonth = _.groupBy(purchases.Data, function(p) {
      var info = p.Activation ? p.Activation : p.Reservation;
      return moment(info.TimeStart).format('MMM YYYY');
    });
    var monthlyBills = _.map(billMonths, function(m) {
      var month = m.format('MMM YYYY');
      var monthlyBill = {
        month: month,
        memberships: [],
        purchases: [],
        sums: {
          memberships: {
            priceInclVAT: 0,
            priceExclVAT: 0,
            priceVAT: 0
          },
          purchases: {
            priceInclVAT: 0,
            priceExclVAT: 0,
            priceVAT: 0,
            durations: 0
          },
          total: {}
        }
      };

      /*
       * Collect purchases and sum for the totals
       */
      _.eachRight(purchasesByMonth[month], function(info) {
        var timeStart;
        var timeEnd;
        if (info.Activation) {
          timeStart = moment(info.Activation.TimeStart);
          timeEnd = moment(info.Activation.TimeEnd);
        } else {
          timeStart = moment(info.Reservation.TimeStart);
          timeEnd = moment(info.Reservation.TimeEnd);
        }

        var duration = moment.duration(timeEnd.diff(timeStart))
                             .asSeconds();
        
        monthlyBill.sums.durations += duration;
        var priceInclVAT = toCents(info.DiscountedTotal);
        var priceExclVAT = toCents(subtractVAT(info.DiscountedTotal));
        var priceVAT = priceInclVAT - priceExclVAT;
        monthlyBill.sums.purchases.priceInclVAT += priceInclVAT;
        monthlyBill.sums.purchases.priceExclVAT += priceExclVAT;
        monthlyBill.sums.purchases.priceVAT += priceVAT;
        monthlyBill.purchases.push({
          MachineName: info.Machine.Name,
          Type: info.Activation ? 'Activation' : 'Reservation',
          TimeStart: timeStart,
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
        priceInclVAT: monthlyBill.sums.purchases.priceInclVAT + monthlyBill.sums.memberships.priceInclVAT,
        priceExclVAT: monthlyBill.sums.purchases.priceExclVAT + monthlyBill.sums.memberships.priceExclVAT,
        priceVAT: monthlyBill.sums.purchases.priceVAT + monthlyBill.sums.memberships.priceVAT
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

const getReservations = [
  ['reservationsStore'],
  (reservationsStore) => {
    return reservationsStore.get('reservations');
  }
];

const getReservationsByDay = [
  getReservations,
  (reservations) => {
    if (reservations) {
      return toImmutable(_.groupBy(reservations.toArray(), (reservation) => {
        return moment(reservation.get('TimeStart')).format('YYYY-MM-DD');
      }));
    }
  }
];

const getActiveReservationsByMachineId = [
  getReservationsByDay,
  (reservationsByDay) => {
    if (reservationsByDay) {
      var m = moment();
      var u = m.unix();
      var rs = reservationsByDay.get(m.format('YYYY-MM-DD'));
      var reservationsByMachineId = {};
      if (rs) {
        _.each(rs.toObject(), (reservation) => {
          var start = moment(reservation.get('start')).unix();
          var end = moment(reservation.get('end')).unix();
          if (start <= u && u <= end) {
            reservationsByMachineId[reservation.get('MachineId')] = reservation;
          }
        });
      }
      return toImmutable(reservationsByMachineId);
    }
  }
];

const getNewReservationTimes = [
  getMachinesById,
  ['reservationsStore'],
  ['reservationRulesStore'],
  getReservationsByDay,
  (machinesById, reservationsStore, reservationRulesStore, reservationsByDay) => {
    var newReservation = reservationsStore.get('create');
    var allMachineIds = [];
    _.each(machinesById.toJS(), function(machine, id) {
      id = parseInt(id);
      allMachineIds.push(id);
    });

    if (newReservation) {
      return newReservation.get('times').map((t) => {
        /* Machine Ids available according to reservation rules */
        var availableIds = reservationRulesStore
          .get('reservationRules')
          .sortBy((rule) => {
            return rule.get('Available');
          })
          .reduce((availableMachineIds, rule) => {
            var applies = true;
            var tm;
            var d;
            if (rule.get('DateStart') && rule.get('TimeStart')) {
              var start = moment(rule.get('DateStart') + ' ' + rule.get('TimeStart'));
              if (start.unix() > t.get('end').unix()) {
                applies = false;
              }
            } else if (rule.get('DateStart')) {
              var dateStart = new Day(rule.get('DateStart'));
              d = new Day(t.get('end').format('YYYY-MM-DD'));
              if (dateStart.toInt() > d.toInt()) {
                applies = false;
              }
            } else if (rule.get('TimeStart')) {
              var timeStart = new Time(rule.get('TimeStart'));
              tm = new Time(t.get('end').format('HH:mm'));
              if (timeStart.toInt() >= tm.toInt()) {
                applies = false;
              }
            }

            if (rule.get('DateEnd') && rule.get('TimeEnd')) {
              var end = moment(rule.get('DateEnd') + ' ' + rule.get('TimeEnd'));
              if (end.unix() < t.get('start').unix()) {
                applies = false;
              }
            } else if (rule.get('DateEnd')) {
              var dateEnd = new Day(rule.get('DateEnd'));
              d = new Day(t.get('end').format('YYYY-MM-DD'));
              if (dateEnd.toInt() < d.toInt()) {
                applies = false;
              }
            } else if (rule.get('TimeEnd')) {
              var timeEnd = new Time(rule.get('TimeEnd'));
              tm = new Time(t.get('start').format('HH:mm'));
              if (timeEnd.toInt() < tm.toInt()) {
                applies = false;
              }
            }

            var anyWeekDaySelected = false;
            _.each(['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'], (weekDay) => {
              if (rule.get(weekDay)) {
                anyWeekDaySelected = true;
              }
            });

            if (!rule.get('DateStart') && !rule.get('TimeStart') && !rule.get('DateEnd') && !rule.get('TimeEnd') && !anyWeekDaySelected) {
              applies = false;
            }

            if (anyWeekDaySelected) {
              switch (reservationsStore.get('create').get('date').isoWeekday()) {
              case 1:
                applies = applies && !!rule.get('Monday');
                break;
              case 2:
                applies = applies && !!rule.get('Tuesday');
                break;
              case 3:
                applies = applies && !!rule.get('Wednesday');
                break;
              case 4:
                applies = applies && !!rule.get('Thursday');
                break;
              case 5:
                applies = applies && !!rule.get('Friday');
                break;
              case 6:
                applies = applies && !!rule.get('Saturday');
                break;
              case 7:
                applies = applies && !!rule.get('Sunday');
                break;
              }
            } else {
              applies = false;
            }

            if (applies && rule.get('Unavailable')) {
              if (rule.get('MachineId')) {
                return _.difference(availableMachineIds, [rule.get('MachineId')]);
              } else {
                return [];
              }
            } else if (applies && rule.get('Available')) {
              if (rule.get('MachineId')) {
                return _.union(availableMachineIds, [rule.get('MachineId')]);
              } else {
                return allMachineIds; 
              }
            } else {
              return availableMachineIds;
            }
          }, allMachineIds);

        /* Consider colliding reservations */
        var day = t.get('start').format('YYYY-MM-DD');
        var rs = reservationsByDay.get(day);
        if (rs) {
          _.each(rs.toArray(), function(reservation) {
            var rStart = moment(reservation.get('TimeStart')).unix();
            var rEnd = moment(reservation.get('TimeEnd')).unix();
            var tStart = moment(t.get('start')).unix();
            var tEnd = moment(t.get('end')).unix();
            var tStartInInterval = tStart >= rStart && tStart < rEnd;
            var tEndInInterval = tEnd > rStart && tEnd <= rEnd;
            if (tStartInInterval || tEndInInterval) {
              availableIds = _.difference(availableIds, [reservation.get('MachineId')]);
            }
          });
        }
        return t.set('availableMachineIds', availableIds);
      });
    } // if (newReservation)
  }
];

const getNewReservationFrom = [
  getNewReservationTimes,
  (reservationTimes) => {
    if (reservationTimes) {
      var selectedTimes = _.filter(reservationTimes.toArray(), function(t) {
        return t.get('selected');
      });
      if (selectedTimes.length > 0) {
        return _.first(selectedTimes).get('start').toDate();
      }
    }
  }
];

const getNewReservationTo = [
  getNewReservationTimes,
  (reservationTimes) => {
    if (reservationTimes) {
      var selectedTimes = _.filter(reservationTimes.toArray(), function(t) {
        return t.get('selected');
      });
      if (selectedTimes.length > 0) {
        return _.last(selectedTimes).get('end').toDate();
      }
    }
  }
];

const getNewReservationPrice = [
  getMachinesById,
  getNewReservation,
  getNewReservationTimes,
  (machinesById, newReservation, times) => {
    if (machinesById && times) {
      var machineId = newReservation.get('machineId');
      var machine = machinesById.get(machineId);
      if (machine) {
        var pricePerSlot = machine.get('ReservationPriceHourly') / 2;
        var slots = _.reduce(times.toJS(), (total, slot) => {
          return total + (slot.selected ? 1 : 0);
        }, 0);
        return slots * pricePerSlot;
      }
    }
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
  getNewReservation, getNewReservationPrice, getNewReservationTimes, getNewReservationFrom, getNewReservationTo, getReservations, getActiveReservationsByMachineId,
  getScrollUpEnabled, getScrollDownEnabled, getScrollPosition
};
