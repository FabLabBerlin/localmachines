var getters = require('../getters');
var MachineActions = require('../actions/MachineActions');
var React = require('react');
var reactor = require('../reactor');
var toastr = require('toastr');


var MaintenanceSwitch = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machinesById: getters.getMachinesById,
      userInfo: getters.getUserInfo
    };
  },

  handleClick(onOrOff) {
    var mid = this.props.machineId;
    MachineActions.setUnderMaintenance({ mid, onOrOff });
  },

  render() {
    const isAdmin = this.state.userInfo.get('UserRole') === 'admin';
    if (isAdmin) {
      const machine = this.state.machinesById[this.props.machineId];
      return (
        <div className="machine-maintenance-switch">
          Under Maintenance: 
          {machine.UnderMaintenance ?
            (
              <i className="fa fa-toggle-on"
                 onClick={this.handleClick.bind(this, 'off')}/>
            ) :
            (
              <i className="fa fa-toggle-off"
                 onClick={this.handleClick.bind(this, 'on')}/>
            )
          }
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default MaintenanceSwitch;
