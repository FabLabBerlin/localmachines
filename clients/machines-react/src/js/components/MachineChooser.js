import React from 'react';
import MachineActions from '../actions/MachineActions';
import OccupiedMachine from './OccupiedMachine';
import BusyMachine from './BusyMachine';
import FreeMachine from './FreeMachine';

/*
 * Multiple button available there !!
 * Has to be connected to activation store
 */
var MachineChooser = React.createClass({

  forceSwitch(onOrOff) {
    if(onOrOff === 'off') {
      MachineActions.adminTurnOffMachine(mid, aid)
    } else if (onOrOff === 'on') {
      MachineActions.adminTurnOnMachine(mid)
    }
    //adminTurnOnMachine(mid)
  },

  endActivation() {
    let aid = this.props.activation.Id
    MachineActions.endActivation(aid);
  },

  startActivation() {
    let mid = this.props.info.Id;
    MachineActions.startActivation(mid)
  },

  /*
   * Not render the component when the props doesn't change;
   */
  shouldComponentUpdate(nextProps) {
    return nextProps.activation.Id !== this.props.activation.Id;
  },

  /*
   * Render a machine component depending of the props
   */
  //<div className="machine-info-content">{this.props.info.Description}</div>
  render() {
    var isAdmin = this.props.user.Role === 'admin';
    return (
      <div className="machine available">
        <div className="container-fluid" >
          <div className="machine-header">
            <div className="machine-title pull-left">{this.props.info.Name}</div>
            <div className="machine-info-btn pull-right">

            </div>
            <div className="clearfix"></div>
          </div>
          <div className="machine-body">
            { this.props.busy ?
              this.props.sameUser ? (
                <BusyMachine
                  activation={this.props.activation}
                  info={this.props.info}
                  isAdmin={isAdmin}
                  func={this.endActivation}
                  force={this.forceSwitch}
                />
            ) : (
            <OccupiedMachine
              activation={this.props.activation}
              info={this.props.info}
              isAdmin={isAdmin}
              func={this.endActivation}
              force={this.forceSwitch}
            />
            ) :
              (
                <FreeMachine
                  info={this.props.info}
                  isAdmin={isAdmin}
                  func={this.startActivation}
                  force={this.forceSwitch}
                />
            )}
          </div>
        </div>
      </div>
    );
  }
});

module.exports = MachineChooser;
