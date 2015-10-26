var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');


var Slot = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservations: getters.getReservations,
      userId: getters.getUid
    };
  },

  render() {
    var now = moment().unix();
    var reserved = false;
    var reservedByUser = false;

    _.each(this.state.reservations.toJS(), r => {
      if (r.MachineId === this.props.machineId) {
        var start = moment(r.TimeStart).unix();
        var end = moment(r.TimeEnd).unix();
        var u = this.props.time.unix();
        if (start <= u && u <= end) {
          reserved = true;
          if (r.UserId === this.state.userId) {
            reservedByUser = true;
          }
        }
      }
    }.bind(this));

    var className = 'slot';
    if (reserved) {
      className += ' reserved';
      if (reservedByUser) {
        className += ' by-user';
      }
    }
    var nn = now - (now % 1800);
    var uu = this.props.time.unix() - (this.props.time.unix() % 1800);
    if (nn === uu) {
      className += ' now';
    }

    return <div className={className}/>;
  }
});


var AvailabilityDisplay = React.createClass({

  render() {
    var times = [];
    var i = 0;
    var n = 2 * 48;

    for (var t = moment().hours(0); i < n; t.add(30, 'm'), i++) {
      times.push(t.clone());
    }

    return (
      <div className="machine-reserv-preview">
        <div className="today">
          <div className="slots">
            {_.map(times.slice(0, n / 2), time => <Slot machineId={this.props.machineId} time={time}/>)}
          </div>
          <div className="label">
            Today
          </div>
        </div>
        <div className="tomorrow">
          <div className="slots">
            {_.map(times.slice(n / 2), time => <Slot machineId={this.props.machineId} time={time}/>)}
          </div>
          <div className="label">
            Tomorrow
          </div>
        </div>
      </div>
    );
  }

});

export default AvailabilityDisplay;
