var moment = require('moment');


class Month {
  // monthId from 0 to 11
  constructor(monthId, year) {
    this._monthId = monthId;
    this._year = year;
    this._firstDay = moment([year, monthId, 1]);
    if (!this._firstDay.isValid()) {
      console.log('this._firstDay is invalid!!!');
    }
    this._lastDay = this._firstDay.clone()
                                  .add(1, 'month')
                                  .subtract(1, 'day');
  }

  static getCurrentMonth() {
    return new Month(moment().month(), moment().year());
  }

  getNextMonth() {
    var t = this.lastDay().add(1, 'day');
    return new Month(t.month(), t.year());
  }

  firstDay() {
    return this._firstDay.clone();
  }

  lastDay() {
    return this._lastDay.clone();
  }

  toString() {
    return this.firstDay().format('MMMM');
  }

  weeks() {
    return (this.lastDay().week() - this.firstDay().week() + 1 + 52) % 52;
  }
}


class Day {
  constructor(yyyymmdd) {
    var yyyy = yyyymmdd.slice(0, 4);
    var mm = yyyymmdd.slice(5, 7);
    var dd = yyyymmdd.slice(8, 10);
    this._years = parseInt(yyyy);
    this._months = parseInt(mm);
    this._days = parseInt(dd);
  }

  toInt() {
    return this._years * 400 + this._months * 31 + this._days;
  }
}


class Time {
  constructor(hhmm) {
    var mm;
    var hh;
    if (hhmm.length === 4) {
      hh = hhmm[0];
      mm = hhmm.slice(2, 4);
    } else {
      hh = hhmm.slice(0, 2);
      mm = hhmm.slice(3, 5);
    }
    this._hours = parseInt(hh, 10);
    this._minutes = parseInt(mm, 10);
  }

  static now() {
    var m = moment();
    return new Time(m.format('HH:mm'));
  }

  isLargerEqual(t) {
    if (this._hours === t._hours) {
      return this._minutes >= t._minutes;
    } else {
      return this._hours >= t._hours;
    }
  }

  toInt() {
    return this._hours * 60 + this._minutes;
  }

  toString() {
    return String(this._hours) + ':' + this._minutes;
  }
}

export default { Day, Month, Time };
