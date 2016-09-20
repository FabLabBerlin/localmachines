var _ = require('lodash');
var $ = require('jquery');
var getters = require('../../getters');
var React = require('react');
var ReactDOM = require('react-dom');
var reactor = require('../../reactor');
var ReservationActions = require('../../actions/ReservationActions');
var Settings = require('../../modules/Settings');


var Time = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      currency: Settings.getters.getCurrency,
      newReservationPrice: getters.getNewReservationPrice
    };
  },

  render() {
    const className = this.props.className;
    const onChange = this.props.onChange;
    const pricing = this.props.pricing;
    const showPrice = this.props.showPrice;
    const t = this.props.t;

    return (
      <div className={className}>
        <label>
          <input
            checked={t.selected}
            type="checkbox"
            onChange={onChange}
          />
          <div className="row">
            <div className="col-md-3">
            </div>
            <div className="col-md-6">
              {t.start.format('HH:mm')} - {t.end.format('HH:mm')}
            </div>
            <div className="col-md-3">
              {showPrice ? (
                <div className="total-price">
                  Total price: {this.state.currency}{(this.state.newReservationPrice).toFixed(2)}
                </div>
              ) : null}
            </div>
          </div>
        </label>
      </div>
    );
  }
});


var TimePicker = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      newReservation: getters.getNewReservation,
      times: getters.getNewReservationTimes
    };
  },

  isRange(times) {
    var lastIsSelected;
    var rangesFound = 0;
    _.each(times, function(t) {
      if (t.selected && !lastIsSelected) {
        rangesFound++;
      }
      lastIsSelected = t.selected;
    });
    return rangesFound === 1;
  },

  render() {
    var machineId = this.state.newReservation.get('machineId');
    var lastIndex = _.reduce(this.state.times.toJS(), (lastIdx, t, i) => {
      return t.selected ? i : lastIdx;
    }, null);

    var containerClassName = 'time-picker ' + this.props.className;

    return (
      <div className={containerClassName}>
        <h3>Select time range</h3>
        <div className="no-select" ref="times">
          {_.map(this.times(), (t, i) => {
            const offset = this.state.times.count() - this.times().length;

            var className = 'time';
            var onChange;

            if (!_.includes(t.availableMachineIds, machineId)) {
              className += ' unavailable';
            } else {
              onChange = this.setTimes.bind(this, i + offset);
              if (t.selected) {
                className += ' selected';
              }
            }

            return <Time key={i}
                         className={className}
                         onChange={onChange}
                         showPrice={(i + offset) === lastIndex}
                         t={t}/>;
          })}
        </div>
        <hr/>
        <div className="pull-right">
          <button className="btn btn-lg btn-info" type="button" onClick={this.previous}>Previous</button>
          <button className="btn btn-lg btn-primary" type="button" onClick={this.submit}>Submit</button>
        </div>
      </div>
    );
  },

  previous() {
    ReservationActions.newReservation.previousStep();
  },

  setTimes(lastClickIndex) {
    var times = this.state.times.toJS();

    var first = null;
    var firstIndex;
    var last = null;
    var lastIndex;
    $(ReactDOM.findDOMNode(this.refs.times)).find('input').each(function(i, el) {

      // Find first checked element
      if (first === null && el.checked) {
        first = el;
        firstIndex = i;
      } 

      // Find last checked element
      if (el.checked) {
        last = el;
        lastIndex = i;
      }

      times[i + this.state.times.count() - this.times().length].selected = false;
    }.bind(this));

    if ((firstIndex && lastIndex) && (firstIndex === lastIndex)) {
      times[firstIndex + this.state.times.count() - this.times().length].selected = true;
    } else {
      var doSelect = false;
      $(ReactDOM.findDOMNode(this.refs.times)).find('input').each(function(i, el) {
        if (parseInt(i) === parseInt(firstIndex)) {
          doSelect = true;
        }

        // Do the selection between first and last checked time item
        if (doSelect) {
          times[i + this.state.times.count() - this.times().length].selected = true;
        } else {
          times[i + this.state.times.count() - this.times().length].selected = false;
        }

        if (parseInt(i) === parseInt(lastIndex)) {
          doSelect = false;
        } 
      }.bind(this));
    }

    ReservationActions.newReservation.setTimes({ times });
  },

  submit() {
    ReservationActions.newReservation.submit();
  },

  times() {
    var times = this.state.times.toJS();
    var machineId = this.state.newReservation.get('machineId');
    var anyTimeAvailable = _.reduce(times, (foundAny, t) => {
      return foundAny || _.includes(t.availableMachineIds, machineId);
    }, false);

    if (!anyTimeAvailable) {
      times = [];
      var i = 0;
      var minHour = 10;
      var maxHour = 19;
      if (this.state.newReservation.get('date').isoWeekday() === 6) {
        minHour = 12;
        maxHour = 18;
      }
      for (var tt = this.state.newReservation.get('date').clone().hours(minHour); i < 2 * (maxHour - minHour); tt.add(30, 'm'), i++) {
        times.push({
          start: tt.clone(),
          end: tt.clone().add(30, 'm'),
          selected: false,
          availableMachineIds: []
        });
      }
    } else {
      var someAvailableAlready = false;
      times = _.reduce(times, (newTimes, ttt) => {
        if (!someAvailableAlready && !_.includes(ttt.availableMachineIds, machineId)) {
          return newTimes;
        } else {
          someAvailableAlready = true;
          newTimes.push(ttt);
          return newTimes;
        }
      }, []);
    }
    return times;
  }
});

export default TimePicker;
