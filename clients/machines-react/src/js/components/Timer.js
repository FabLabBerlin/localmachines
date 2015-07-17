import React from 'react';

var Timer = React.createClass({

  getInitialState() {
    var now = new Date().getTime();
    var startTime = new Date(this.props.time).getTime();
    var initialTime = ( now - startTime )/ 1000;
    return {secondsElpased: initialTime};
  },

  tick() {
    this.setState({
      secondsElapsed: this.state.secondsElpased + 1
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
      <div>{this.state.secondsElapsed}</div>
    );
  }
});

module.exports = Timer;
