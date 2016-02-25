var React = require('react');
var getters = require('../../../getters');
var reactor = require('../../../reactor');

/*
 * Methode to int to be display in hour format
 * Use in timer to display the time you spend
 */
function toHHMMSS(d) {
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
    return {
      secondsElapsed: this.props.activation.Quantity * this.multiplier()
    };
  },

  multiplier() {
    switch (this.props.activation.PriceUnit) {
    case 'day':
      return 86400;
    case 'hour':
      return 3600;
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
      secondsElapsed: this.state.secondsElapsed + 1
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
    const uid = reactor.evaluateToJS(getters.getUid);
    const sameUser = uid === this.props.activation.UserId;
    const machinesById = reactor.evaluateToJS(getters.getMachinesById);
    const machine = machinesById[this.props.activation.MachineId];
    var inGracePeriod = machine && machine.GracePeriod && 
                        duration < machine.GracePeriod;
    var duration = this.state.secondsElapsed;
    if (inGracePeriod) {
      duration = machine.GracePeriod - duration;
    } else {
      duration = duration - machine.GracePeriod;
    }
    return (
      <div className="machine-time-value">
        {sameUser ? (
          <div className="machine-time-label">
            {inGracePeriod ?
                'Grace Period'
              : 'Usage time'}
          </div>
        ) : null}
        {toHHMMSS(duration)}
      </div>
    );
  }
});

export default Timer;
