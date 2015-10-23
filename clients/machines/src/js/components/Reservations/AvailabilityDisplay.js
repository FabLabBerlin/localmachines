var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');


var Slot = React.createClass({
  render() {
    var className = 'slot ';
    if (this.props.busy) {
      className += 'busy';
    }
    return <div className={className}/>;
  }
});


var AvailabilityDisplay = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservations: getters.getReservations
    };
  },

  render() {
    var slots = [];
    var i = 0;

    for (var t = moment().hours(0); i < 2 * 48; t.add(30, 'm'), i++) {
      var reserved = false;
      _.each(this.state.reservations.toJS(), r => {
          if (r.MachineId === this.props.machineId) {
          var start = moment(r.TimeStart).unix();
          var end = moment(r.TimeEnd).unix();
          var u = t.unix();
          if (start <= u && u <= end) {
            reserved = true;
          }
        }
      });
      slots.push(reserved);
    }

    return (
      <div className="machine-reserv-preview">
        {_.map(slots, busy => <Slot busy={busy}/>)}
      </div>
    );
  }

});

export default AvailabilityDisplay;
