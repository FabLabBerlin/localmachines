var getters = require('../../../getters');
var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('toastr');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';


var MaintenanceSwitch = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machinesById: getters.getMachinesById,
      user: getters.getUser
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
      callback: function(confirmed) {
        if (confirmed) {
          MachineActions.setUnderMaintenance({ mid, onOrOff });
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }.bind(this)
    });

  },

  render() {
    const isAdmin = this.state.user.get('UserRole') === 'admin';
    if (isAdmin) {
      const machine = this.state.machinesById.get(this.props.machineId);
      return (
        <div className="machine-maintenance-switch">
          {machine.get('UnderMaintenance') ? (
            <a 
              className="danger" 
              href="#" 
              onClick={this.handleClick.bind(this, 'off')}>
              <i className="fa fa-toggle-on"></i>
            </a>
          ) : (
            <a 
              className="primary" 
              href="#" 
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
