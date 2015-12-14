var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var TutoringActions = require('../../actions/TutoringActions');

var Tutoring = React.createClass({

  getInitialState() {
    return {
      secondsElapsed: 0
    };
  },

  componentDidMount() {
    this.interval = setInterval(this.tick, 1000);
  },

  componentWillUnmount() {
    clearInterval(this.interval);
  },

  startTimer() {
    this.props.tutoring.TimerTimeStart = moment().format('YYYY-MM-DD HH:mm:ss');
    TutoringActions.startTutoring(this.props.tutoring.Id);
  },

  stopTimer() {
    TutoringActions.stopTutoring(this.props.tutoring.Id);
  },

  tick() {
    this.setState({
      secondsElapsed: this.state.secondsElapsed + 1
    });
  },

  render() {
    var start = moment(this.props.tutoring.TimeStart);
    var end = moment(this.props.tutoring.TimeEnd);
    var machineUsers = reactor.evaluateToJS(getters.getMachineUsers);

    var currentTimerDuration;
    if (this.props.tutoring.PriceUnit === 'day') {
      currentTimerDuration = moment.duration(this.props.tutoring.Quantity, 'days');
    } else if (this.props.tutoring.PriceUnit === 'hour') {
      currentTimerDuration = moment.duration(this.props.tutoring.Quantity, 'hours');
    } else if (this.props.tutoring.PriceUnit === 'minute') {
      currentTimerDuration = moment.duration(this.props.tutoring.Quantity, 'minutes');
    }

    var duration;
    var user;

    if (end) {
      end = moment(end);
      if (end.unix() < 0) {
        end = null;
      }
    }

    if (machineUsers) {
      user = machineUsers[this.props.tutoring.UserId];
    }

    if (this.props.tutoring.Running) {
      var lastTimerStart = moment(this.props.tutoring.TimerTimeStart);
      var now = moment();
      var elapsed = moment.duration(now.diff(lastTimerStart));

      duration = currentTimerDuration.add(elapsed).format('h[h] m[m] s[s]');
    } else if (start && end) {
      duration = currentTimerDuration.format('h[h] m[m] s[s]');
    }

    return (
      <div className="tutoring-item">
        <div className="container-fluid">
          <div className="row">
            <div className="col-xs-8">
              <div className="row">
                <div className="col-xs-6">
                  <div><b>User</b></div>
                  <div>{user && (user.FirstName + ' ' + user.LastName)}</div>
                </div>
                <div className="col-xs-6">
                  <div><b>Timer</b></div>
                  <div>{duration}</div>
                </div>
              </div>
              <div className="row">
                <div className="col-xs-6">
                  <div><b>Time from</b></div>
                  <div>{start && start.format('D MMM YYYY HH:mm')}</div>
                </div>
                <div className="col-xs-6">
                  <div><b>Time to</b></div>
                  <div>{end && end.format('D MMM YYYY HH:mm')}</div>
                </div>
              </div>
            </div>
            <div className="col-xs-4">

              {this.props.tutoring.Running ? (
                <button 
                  className="btn btn-danger btn-lg btn-block"
                  onClick={this.stopTimer}>
                  Stop
                </button>
              ) : (
                <button 
                  className="btn btn-primary btn-lg btn-block"
                  onClick={this.startTimer}>
                  Start
                </button>
              )}

            </div>
          </div>
        </div>
      </div>
    );
  }
});

export default Tutoring;
