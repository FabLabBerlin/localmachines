var $ = require('jquery');
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');


var TimePicker = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      times: getters.getNewReservationTimes
    };
  },

  render() {
    return (
      <div>
        <h3 className="h3">Select time range</h3>
        <div className="no-select" ref="times">
          {_.map(this.state.times.toJS(), (t, i) => {
            var className = 'time-picker-time';
            if (t.selected) {
              className += ' selected';
            }
            return (
              <div key={i} className={className}>
                <label>
                  <input
                    type="checkbox"
                    onChange={this.setTimes}
                  />
                  {t.start.format('HH:mm')} - {t.end.format('HH:mm')}
                </label>
              </div>
            );
          })}
        </div>
        <button className="btn btn-lg btn-info" type="button" onClick={this.previous}>Previous</button>
        <button className="btn btn-lg btn-primary" type="button" onClick={this.submit}>Submit</button>
      </div>
    );
  },

  previous() {
    ReservationsActions.previousStep();
  },

  setTimes() {
    var times = this.state.times.toJS();
    $(this.refs.times.getDOMNode()).find('input').each(function(i, el) {
      times[i].selected = el.checked;
    });
    ReservationsActions.createSetTimes({ times });
  },

  submit() {
    ReservationsActions.createSubmit();
  }
});

export default TimePicker;
