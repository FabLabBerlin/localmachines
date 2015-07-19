import React from 'react';

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
      <div>{this.state.secondsElapsed}</div>
    );
  }
});

module.exports = Timer;
