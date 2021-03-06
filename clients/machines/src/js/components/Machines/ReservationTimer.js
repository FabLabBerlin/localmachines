import React from 'react';
import getters from '../../getters';
import Machines from '../../modules/Machines';
import moment from 'moment';
import reactor from '../../reactor';

/*
 * Methode to int to be display in hour format
 * Use in timer to display the time you spend
 */
function toHHMMSS(d) {
  if (d < 0) {
    d *= -1;
  }
  var h = Math.floor(d / 3600);
  var m = Math.floor(d % 3600 / 60);
  var s = Math.floor(d % 3600 % 60);

  return ((h > 0 ? h + ':' + (m < 10 ? '0' : '') : '') + m + ':' + (s < 10 ? '0' : '') + s);
}


/*
 * Timer who will display the time you use the machine
 */
var Timer = React.createClass({

  /*
   * Initial State
   * Take the Quantity from activation json
   */
  getInitialState() {
    this.startShow = moment().unix();
    return {
      timeNow: moment().unix(),
      timeStart: moment(this.props.reservation.TimeStart).unix()
    };
  },

  multiplier() {
    switch (this.props.activation.PriceUnit) {
    case 'day':
      return 86400;
    case 'hour':
      return 3600;
    case '30 minutes':
      return 1800;
    case 'minute':
      return 60;
    case 'second':
      return 1;
    default:
      return 0;
    }
  },

  /*
   * Function which will increase the state each seconds
   * It will rerender the component
   */
  tick() {
    this.setState({
      timeNow: moment().unix()
    });
  },

  /*
   * Destructor
   * Clear the interval when the component
   */
  componentWillUnmount() {
    clearInterval(this.interval);
  },

  /*
   * Call once when component in mount on the DOM
   * Start interval and call tick each seconds
   */
  componentDidMount() {
    this.interval = setInterval(this.tick, 1000);
  },

  /*
   * Render the timer
   */
  render() {
    var duration = this.state.timeStart - this.state.timeNow;

    return (
      <div className="machine-time-value">
        {toHHMMSS(duration)}
      </div>
    );
  }
});

export default Timer;
