import React from 'react';
import OccupiedMachine from './OccupiedMachine';
import BusyMachine from './BusyMachine';
import FreeMachine from './FreeMachine';

/*
 * Multiple button available there !!
 * Has to be connected to activation store
 */
var MachineChooser = React.createClass({
  render() {
    return (
      <div>
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
          activation={this.props.activation}
          info={this.props.info}
        />
        )}
      </div>
    );
  }
});

module.exports = MachineChooser;
