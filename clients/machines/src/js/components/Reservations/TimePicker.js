var $ = require('jquery');
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');


var TimePicker = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machinesById: getters.getMachinesById,
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
    var machine = this.state.machinesById.get(machineId);
    var pricePerSlot = machine.get('ReservationPriceHourly') / 2;
    var lastIndex = _.reduce(this.state.times.toJS(), (lastIdx, t, i) => {
      if (t.selected) {
        return i;
      } else {
        return lastIdx;
      }
    }, null);

    return (
      <div className={this.props.className}>
        <h3 className="h3">Select time range</h3>
        <div className="no-select" ref="times">
          {_.map(this.state.times.toJS(), (t, i) => {
            var className = 'time-picker-time';
            var onChange;
            var pricingInfo;

            if (!_.includes(t.availableMachineIds, machineId)) {
              className += ' unavailable';
            } else {
              onChange = this.setTimes.bind(this, i);
              if (t.selected) {
                className += ' selected';
              }
            }

            if (i === lastIndex) {
              var slots = _.reduce(this.state.times.toJS(), (total, slot) => {
                return total + (slot.selected ? 1 : 0);
              }, 0);
              pricingInfo = (
                <div>
                  Total price: {(slots * pricePerSlot).toFixed(2)} €
                </div>
              );
            }

            return (
              <div key={i} className={className}>
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
                      {pricingInfo}
                    </div>
                  </div>
                </label>
              </div>
            );
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
    ReservationsActions.previousStep();
  },

  setTimes(lastClickIndex) {
    var times = this.state.times.toJS();
    $(this.refs.times.getDOMNode()).find('input').each(function(i, el) {
      times[i].selected = el.checked;
    });
    if (!this.isRange(times)) {
      _.each(times, function(t, i) {
        t.selected = lastClickIndex === i;
      });
    }
    ReservationsActions.createSetTimes({ times });
  },

  submit() {
    ReservationsActions.createSubmit();
  }
});

export default TimePicker;
