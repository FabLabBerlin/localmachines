var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var TutoringActions = require('../../actions/TutoringActions');


var Tutoring = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: getters.getMachineUsers
    };
  },

  startTimer() {
    TutoringActions.startTutoring(this.props.tutoring.Id);
  },

  stopTimer() {
    TutoringActions.stopTutoring(this.props.tutoring.Id);
  },

  render() {
    var start = moment(this.props.tutoring.TimeStart);
    var end = this.props.tutoring.TimeEnd || this.props.tutoring.TimeEndPlanned;
    var duration;
    var user;

    if (end) {
      end = moment(end);
      if (end.unix() < 0) {
        end = null;
      }
    }

    if (this.state.machineUsers) {
      user = this.state.machineUsers.get(this.props.tutoring.UserId);
    }

    if (this.props.tutoring.Running) {
      duration = moment().subtract(start);
      duration = duration.format('HH:mm:ss');
    } else if (start && end) {
      duration = end.clone().subtract(start).format('HH:mm:ss');
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
                  <div>{start && start.format('YYYY-MM-DD HH:mm')}</div>
                </div>
                <div className="col-xs-6">
                  <div><b>Time to</b></div>
                  <div>{end && end.format('YYYY-MM-DD HH:mm')}</div>
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
