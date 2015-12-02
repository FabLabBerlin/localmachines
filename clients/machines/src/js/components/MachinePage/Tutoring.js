var moment = require('moment');
var React = require('react');
var TutoringActions = require('../../actions/TutoringActions');


var Tutoring = React.createClass({

  startTimer() {
    TutoringActions.startTutoring(this.props.tutoring.Id);
  },

  stopTimer() {
    TutoringActions.stopTutoring(this.props.tutoring.Id);
  },

  render() {
    console.log('tutoring: ', this.props.tutoring);
    var start = moment(this.props.tutoring.TimeStart);
    var end = moment(this.props.tutoring.TimeEndActual || this.props.tutoring.TimeEnd);
    var duration;

    if (this.props.tutoring.Running) {
      duration = moment().subtract(start);
      duration = duration.format('HH:mm:ss');
    }

    return (
      <div className="tutoring-item">
        <div className="container-fluid">
          <div className="row">
            <div className="col-xs-8">
              <div className="row">
                <div className="col-xs-6">
                  <div><b>User</b></div>
                  <div>Millumin Atari</div>
                </div>
                <div className="col-xs-6">
                  <div><b>Timer</b></div>
                  <div>{duration}</div>
                </div>
              </div>
              <div className="row">
                <div className="col-xs-6">
                  <div><b>Time from</b></div>
                  <div>{start.format('YYYY-MM-DD HH:mm')}</div>
                </div>
                <div className="col-xs-6">
                  <div><b>Time to</b></div>
                  <div>{end.format('YYYY-MM-DD HH:mm')}</div>
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
