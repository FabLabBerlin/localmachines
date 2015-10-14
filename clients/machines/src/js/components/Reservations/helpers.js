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
    return this.lastDay().week() - this.firstDay().week() + 1;
  }
}

export default { Month };
