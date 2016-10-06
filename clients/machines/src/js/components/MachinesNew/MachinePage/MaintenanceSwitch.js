var $ = require('jquery');
var getters = require('../../../getters');
var LocationGetters = require('../../../modules/Location/getters');
var MachineActions = require('../../../actions/MachineActions');
var Machines = require('../../../modules/Machines');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('toastr');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

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
  render() {
    if (this.props.machine && !this.props.machine.get('UnderMaintenance')) {
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
  render() {
    if (this.props.machine && this.props.machine.get('UnderMaintenance')) {
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
