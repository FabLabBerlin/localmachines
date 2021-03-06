var $ = require('jquery');
import getters from '../../../getters';
import LocationGetters from '../../../modules/Location/getters';
import MachineActions from '../../../actions/MachineActions';
import Machines from '../../../modules/Machines';
import React from 'react';
import reactor from '../../../reactor';
import toastr from 'toastr';

// https://github.com/HubSpot/vex/issues/72
import vex from 'vex-js';
import VexDialog from 'vex-js/js/vex.dialog.js';

vex.defaultOptions.className = 'vex-theme-custom';

function handleClick(machine, onOrOff) {
  VexDialog.buttons.YES.text = 'Yes';
  VexDialog.buttons.NO.text = 'No';

  VexDialog.confirm({
    message: 'Do you really want to set "Under Maintenance" to ' + onOrOff + ' for <b>' +
      machine.get('Name') + '</b>?',
    callback(confirmed) {
      if (confirmed) {
        const mid = machine.get('Id');
        MachineActions.setUnderMaintenance({ mid, onOrOff });
      }
      $('.vex').remove();
      $('body').removeClass('vex-open');
    }
  });
}


var On = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isStaff: LocationGetters.getIsStaff
    };
  },

  render() {
    if (this.props.machine && !this.props.machine.get('UnderMaintenance')
        && this.state.isStaff) {
      return (
        <div className="m-maintenance-switch m-maintenance-action"
             onClick={handleClick.bind(this, this.props.machine, 'on')}>
          Turn on maintenance mode (currently off)
        </div>
      );
    } else {
      return <div/>;
    }
  }
});


var Off = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isStaff: LocationGetters.getIsStaff
    };
  },

  render() {
    if (this.props.machine && this.props.machine.get('UnderMaintenance')
        && this.state.isStaff) {
      return (
        <div className="m-maintenance-switch m-maintenance-action"
             onClick={handleClick.bind(this, this.props.machine, 'off')}>
          Turn off maintenance mode
        </div>
      );
    } else {
      return <div/>;
    }
  }
});


export default {
  On,
  Off
};
