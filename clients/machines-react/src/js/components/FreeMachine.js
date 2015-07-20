import React from 'react';
import MachineActions from '../actions/MachineActions';

/*
 * Div displayed the machine is free
 * Can activate an activation
 */
var FreeMachine = React.createClass ({

  /*
   * Force the switch to turn on
   */
  handleForceSwitchOn() {
    this.props.force('on');
  },

  /*
   * Force the switch to trun off
   */
  handleForceSwitchOff() {
    this.props.force('off');
  },

  /*
   * Try to activate the machine
   */
  startActivation() {
    this.props.func();
  },

  /*
   * Render stuff
   * TODO: real commentaries
   */
  render() {
    console.log('props:', this.props);
    var imageUrl;
    if (this.props.info && this.props.info.machine && this.props.info.machine.Image) {
      imageUrl = this.props.info.machine.Image;
    } else {
      imageUrl = '/machines/assets/img/img-machine-placeholder.svg';
    }
    return (
      <div className="container-fluid">
        <div className="col-xs-6">
          {this.props.activation}
          <div className="machine-action-info">
            <img className="machine-image" 
                 src={imageUrl}/>
          </div>
        </div>
        <div className="col-xs-6">
          <button
            className="btn btn-lg btn-primary btn-block"
            onClick={this.startActivation}
            >Start </button>
          { this.props.isAdmin ? (
            <div className="pull-right" >
              <label>Force Switch</label>
              <button 
                onClick={this.handleForceSwitchOn}
                className="btn btn-lg btn-primary" >On</button>
              <button 
                onClick={this.handleForceSwitchOff}
                className="btn btn-lg btn-danger" >Off</button>
            </div>
          ):('') }
        </div>
      </div>
    );
  }
});

module.exports = FreeMachine;
