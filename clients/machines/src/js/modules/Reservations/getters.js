var Machine = require('../Machine');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


/*
 * Reservations related getters
 */
function slotAvailabilities({date, machinesById, reservationsStore, reservationRulesStore, reservationsByDay}) {
  var allMachineIds = [];
  _.each(machinesById.toJS(), function(machine, id) {
    id = parseInt(id);
    allMachineIds.push(id);
  });

  var timesOfDay = [];
  var i = 0;
  for (var tt = date.clone().hours(0); i < 2 * 24; tt.add(30, 'm'), i++) {
    timesOfDay.push({
      start: tt.clone(),
      end: tt.clone().add(30, 'm'),
      selected: false
    });
  }
  timesOfDay = toImmutable(timesOfDay);

  return timesOfDay.map((t) => {
    /* Machine Ids available according to reservation rules */
    var availableIds = (reservationRulesStore
      .get('reservationRules') || toImmutable({}))
      .sortBy((rule) => {
        /* Unavailability overrides Availability */
        return !rule.get('Available');
      })
      .reduce((availableMachineIds, rule) => {
        var applies = true;
        var tm;
        var d;
        if (rule.get('DateStart')) {
          var dateStart = new Day(rule.get('DateStart'));
          d = new Day(t.get('end').format('YYYY-MM-DD'));
          if (dateStart.toInt() > d.toInt()) {
            applies = false;
          }
        }
        if (rule.get('TimeStart')) {
          var timeStart = new Time(rule.get('TimeStart'));
          tm = new Time(t.get('end').format('HH:mm'));
          if (timeStart.toInt() >= tm.toInt()) {
            applies = false;
          }
        }

        if (rule.get('DateEnd')) {
          var dateEnd = new Day(rule.get('DateEnd'));
          d = new Day(t.get('end').format('YYYY-MM-DD'));
          if (dateEnd.toInt() < d.toInt()) {
            applies = false;
          }
        }
        if (rule.get('TimeEnd')) {
          var timeEnd = new Time(rule.get('TimeEnd'));
          tm = new Time(t.get('start').format('HH:mm'));
          if (timeEnd.toInt() <= tm.toInt()) {
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
          switch (date.isoWeekday()) {
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
      }, []);

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
}

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
          var start = moment(reservation.get('TimeStart')).unix();
          var end = moment(reservation.get('TimeEnd')).unix();
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
  Machine.getters.getMachinesById,
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

    if (newReservation && newReservation.get('date')) {
      if (newReservation.get('times').count() > 0) {
        return newReservation.get('times');
      } else {
        var date = newReservation.get('date');
        return slotAvailabilities({
          date,
          machinesById,
          reservationsStore,
          reservationRulesStore,
          reservationsByDay
        });
      }
    }
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
  Machine.getters.getMachinesById,
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

// getSlotAvailabilities48h returns a map...
// machineId => [reservations in the next 48h]
const getSlotAvailabilities48h = [
  Machine.getters.getMachinesById,
  ['reservationsStore'],
  getReservations,
  (machinesById, reservationsStore, reservations) => {
    //console.log('fgsdfgdfgdf reservations=', reservations);
    var todayStart = moment().hours(0);
    var todayEnd = todayStart.clone().add(1, 'day');
    var tomorrowStart = todayEnd.clone();
    var tomorrowEnd = tomorrowStart.clone().add(1, 'day');
    todayStart = todayStart.unix();
    todayEnd = todayEnd.unix();
    tomorrowStart = tomorrowStart.unix();
    tomorrowEnd = tomorrowEnd.unix();

    if (!reservations) {
      return toImmutable({});
    }

    return reservations.groupBy(reservation => {
      return reservation.get('MachineId');
    }).map((rs) => {
      var tmp = {
        today: rs.filter(r => {
          var start = moment(r.get('TimeStart')).unix();
          var end = moment(r.get('TimeEnd')).unix();
          return start >= todayStart && end <= todayEnd;
        }),
        tomorrow: rs.filter(r => {
          var start = moment(r.get('TimeStart')).unix();
          var end = moment(r.get('TimeEnd')).unix();
          return start >= tomorrowStart && end <= tomorrowEnd;
        })
      };
      return toImmutable(tmp);
    });
  }
];

export default {
  getNewReservation, getNewReservationPrice, getNewReservationTimes, getNewReservationFrom, getNewReservationTo, getReservations, getReservationsByDay, getActiveReservationsByMachineId, getSlotAvailabilities48h
};
