var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');


var Slot = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservations: getters.getReservations
    };
  },

  render() {
    var reserved = false;
    _.each(this.state.reservations.toJS(), r => {
        if (r.MachineId === this.props.machineId) {
        var start = moment(r.TimeStart).unix();
        var end = moment(r.TimeEnd).unix();
        var u = this.props.time.unix();
        if (start <= u && u <= end) {
          reserved = true;
        }
      }
    }.bind(this));
    var className = 'slot ';
    if (reserved) {
      className += 'busy';
    }
    return <div className={className}/>;
  }
});


var AvailabilityDisplay = React.createClass({

  render() {
    var times = [];
    var i = 0;

    for (var t = moment().hours(0); i < 2 * 48; t.add(30, 'm'), i++) {
      times.push(t.clone());
    }

    return (
      <div className="machine-reserv-preview">
        {_.map(times, time => <Slot machineId={this.props.machineId} time={time}/>)}
      </div>
    );
  }

});

export default AvailabilityDisplay;
