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

  /*
   * Render a machine component depending of the props
   */
  //<div className="machine-info-content">{this.props.info.Description}</div>
  render() {
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
                />
            ) : (
            <OccupiedMachine
              activation={this.props.activation}
              info={this.props.info}
              uid={this.props.uid}
            />
            ) :
              (
                <FreeMachine
                  info={this.props.info}
                />
            )}
          </div>
        </div>
      </div>
    );
  }
});

module.exports = MachineChooser;
