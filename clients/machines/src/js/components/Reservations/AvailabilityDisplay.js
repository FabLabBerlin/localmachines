var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');


function isNow(t) {
  var now = moment().unix();
  var nn = now - (now % 1800);
  var uu = t.unix() - (t.unix() % 1800);
  return nn === uu;
}


var Slot = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      userId: getters.getUid
    };
  },

  render() {
    var reserved = true;
    var reservedByUser = this.props.reservation && this.props.reservation.UserId === this.state.userId;

    var className = 'slot';
    if (reserved) {
      className += ' reserved';
      if (reservedByUser) {
        className += ' by-user';
      }
    }
    if (isNow(this.props.time)) {
      className += ' now';
    }

    var style = {
      marginLeft: String(1.04166666 * this.props.position) + '%'
    };

    return <div className={className}
                style={style}/>;
  }
});


var AvailabilityDisplay = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservations: getters.getReservations,
      userId: getters.getUid
    };
  },

  getReservation(t) {
    var reservation;
    _.each(this.state.reservations.toJS(), r => {
      if (r.MachineId === this.props.machineId) {
        var start = moment(r.TimeStart).unix();
        var end = moment(r.TimeEnd).unix();
        var u = t.unix();
        if (start <= u && u <= end) {
          reservation = r;
          return false;
        }
      }
    }.bind(this));
    return reservation;
  },

  render() {
    if (!this.state.reservations) {
      return <div/>;
    }

    var times = [];
    var i = 0;
    var n = 2 * 48;
    var reservation;

    for (var t = moment().hours(0); i < n; t.add(30, 'm'), i++) {
      times.push(t.clone());
    }

    var key = 1;

    return (
      <div className="machine-reserv-preview">
        <div className="today">
          <div className="slots">
            {_.map(times.slice(0, n / 2), (time, j) => {
              reservation = this.getReservation(time);
              if (isNow(time) || reservation) {
                return <Slot key={key++}
                             machineId={this.props.machineId}
                             position={j}
                             reservation={reservation}
                             time={time}/>;
              }
            })}
          </div>
          <div className="label">
            Today
          </div>
        </div>
        <div className="tomorrow">
          <div className="slots">
            {_.map(times.slice(n / 2), (time, j) => {
              reservation = this.getReservation(time)
              if (isNow(time) || reservation) {
                return <Slot key={key++}
                             machineId={this.props.machineId}
                             position={j}
                             reservation={reservation}
                             time={time}/>;
              }
            })}
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
