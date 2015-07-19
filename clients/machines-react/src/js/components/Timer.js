import React from 'react';

String.prototype.toHHMMSS = function() {
  var d = parseInt(this,10);
  var h = Math.floor(d / 3600);
  var m = Math.floor(d % 3600 / 60);
  var s = Math.floor(d % 3600 % 60);
  return ((h > 0 ? h + ':' + (m < 10 ? "0" : "") : "") + m + ":" + (s < 10 ? "0" : "") + s);
}

var Timer = React.createClass({

  getInitialState() {
    return {secondsElapsed: this.props.time};
  },

  tick() {
    this.setState({
      secondsElapsed: this.state.secondsElapsed + 1
    });
  },

  componentWillUnmount: function() {
    clearInterval(this.interval);
  },

  componentDidMount() {
    this.interval = setInterval(this.tick, 1000);
  },

  render() {
    return (
      <div className="machine-time-value">{this.state.secondsElapsed.toString().toHHMMSS()}</div>
    );
  }
});

module.exports = Timer;
