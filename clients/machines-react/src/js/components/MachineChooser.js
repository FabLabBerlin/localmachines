import React from 'react';
import OccupiedMachine from './OccupiedMachine';
import BusyMachine from './BusyMachine';
import FreeMachine from './FreeMachine';

/*
 * Multiple button available there !!
 * Has to be connected to activation store
 */
var MachineChooser = React.createClass({

  toggleInfo() {
    alert('toggleInfo()');
  },

  /*
   * Not render the component when the props doesn't change;
   */
  shouldComponentUpdate(nextProps) {
    return nextProps.activation.Id !== this.props.activation.Id;
  },

  render() {
    return (
      <div className="machine available">
        <div className="machine-header">
          <div className="machine-title pull-left">{this.props.info.Name}</div>
          <div className="machine-info-btn pull-right">

            <a className="machine-info-toggle" onClick={this.toggleInfo}>
              <span className="glyphicon glyphicon-info-sign" ng-class="{hidden: infoVisible}"></span>
              <span className="glyphicon glyphicon-remove-circle" ng-class="{hidden: !infoVisible}"></span>
            </a>

          </div>
          <div className="clearfix"></div>
          <div className="machine-info-content">{this.props.info.Description}</div>
        </div>
        <div className="machine-body">
          { this.props.busy ?
            this.props.sameUser ? (<BusyMachine
              activation={this.props.activation}
              info={this.props.info}
            />
            ) : (
            <OccupiedMachine
              activation={this.props.activation}
              info={this.props.info}
              user={this.props.user}
            />
            ) :(
            <FreeMachine
              info={this.props.info}
            />
            )}
        </div>
      </div>
    );
  }
});

module.exports = MachineChooser;
