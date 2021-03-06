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


var MaintenanceSwitch = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machinesById: Machines.getters.getMachinesById,
      user: getters.getUser,
      isStaff: LocationGetters.getIsStaff
    };
  },

  handleClick(onOrOff) {
    const mid = this.props.machineId;
    const machine = this.state.machinesById.get(mid);
    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to set "Under Maintenance" to ' + onOrOff + ' for <b>' +
        machine.get('Name') + '</b>?',
      callback(confirmed) {
        if (confirmed) {
          MachineActions.setUnderMaintenance({ mid, onOrOff });
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }
    });

  },

  render() {
    if (this.state.isStaff) {
      const machine = this.state.machinesById.get(this.props.machineId);
      return (
        <div className="machine-maintenance-switch">
          {machine.get('UnderMaintenance') ? (
            <a 
              className="danger" 
              href="javascript:void(0)" 
              onClick={this.handleClick.bind(this, 'off')}>
              <i className="fa fa-toggle-on"></i>
            </a>
          ) : (
            <a 
              className="primary" 
              href="javascript:void(0)" 
              onClick={this.handleClick.bind(this, 'on')}>
              <i className="fa fa-toggle-off"></i>
            </a>
          )}Maintenance Mode ({
            machine.get('UnderMaintenance') ? 
            'currently on' : 
            'currently off'
          })
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default MaintenanceSwitch;
