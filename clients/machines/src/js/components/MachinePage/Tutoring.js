var React = require('react');

var Tutoring = React.createClass({

  getInitialState() {
    return {
      timerRunning: false
    };
  },

  startTimer() {
    this.setState({timerRunning: true});
  },

  stopTimer() {
    this.setState({timerRunning: false});
  },

  render() {
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
                  <div>0h 32m</div>
                </div>
              </div>
              <div className="row">
                <div className="col-xs-6">
                  <div><b>Time from</b></div>
                  <div>Today 15:00</div>
                </div>
                <div className="col-xs-6">
                  <div><b>Time to</b></div>
                  <div>Today 17:00</div>
                </div>
              </div>
            </div>
            <div className="col-xs-4">

              {this.state.timerRunning ? (
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
