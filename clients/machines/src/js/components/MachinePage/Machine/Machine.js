var _ = require('lodash');
var $ = require('jquery');
var AvailabilityDisplay = require('../../Reservations/AvailabilityDisplay');
var FeedbackDialogs = require('../../Feedback/FeedbackDialogs');
var getters = require('../../../getters');
var React = require('react');
var LoginActions = require('../../../actions/LoginActions');
var MachineActions = require('../../../actions/MachineActions');
var OccupiedMachine = require('./Occupied');
var BusyMachine = require('./Busy');
var FreeMachine = require('./Free');
var MaintenanceSwitch = require('./MaintenanceSwitch');
var reactor = require('../../../reactor');
var RepairButton = require('../../Feedback/RepairButton');
var ReservedMachine = require('./Reserved');
var UnavailableMachine = require('./Unavailable');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';

/*
 * Multiple machine div available here
 * The component choose which div fit for the situation
 * The choice is made looking at the props
 */
var MachineChooser = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isStaff: getters.getIsStaff,
      machines: getters.getMachines,
      reservationsByMachineId: getters.getActiveReservationsByMachineId
    };
  },

  /*
   * Function pass by props to children
   * End an activation
   */
  endActivation() {
    let aid = this.props.activation.Id;

    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to stop the activation for <b>' +
        this.props.machine.Name + '</b>?',
      callback(confirmed) {
        if (confirmed) {
          MachineActions.endActivation(aid, function() {
            //FeedbackDialogs.checkSatisfaction(aid);
          });
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }
    });

    LoginActions.keepAlive();
  },

  /*
   * Function pass by props to children
   * Start an activation
   */
  startActivation() {
    let mid = this.props.machine.Id;
    MachineActions.startActivation(mid);
    LoginActions.keepAlive();
  },

  /*
   * Render component
   * Can choose what component will be display depending on the props
   * @busy + @sameUser => BusyMachine
   * @busy => OccupiedMachine
   * @nothing => FreeMachine
   */
  render() {
    let isStaff = this.state.isStaff;
    var reservation;
    if (this.state.reservationsByMachineId) {
      reservation = this.state.reservationsByMachineId.toObject()[this.props.machine.Id];
    }
    var machineBody;
    if (reservation && !this.props.busy && !reservation.get('ReservationDisabled') && !reservation.get('Cancelled')) {
      machineBody = (
        <ReservedMachine
          activation={this.props.activation}
          busy={this.props.busy}
          machine={this.props.machine}
          isStaff={isStaff}
          endActivation={this.endActivation}
          startActivation={this.startActivation}
          force={this.forceSwitch}
          reservation={reservation}
          user={this.props.user}
        />
      );
    } else if (this.props.machine.UnderMaintenance) {
      machineBody = (
        <UnavailableMachine
          activation={this.props.activation}
          busy={this.props.busy}
          machine={this.props.machine}
          isStaff={isStaff}
          endActivation={this.endActivation}
          startActivation={this.startActivation}
          force={this.forceSwitch}
        />
      );
    } else {
      if (this.props.busy) {
        if (this.props.sameUser) {
          machineBody = (
            <BusyMachine
              activation={this.props.activation}
              machine={this.props.machine}
              isStaff={isStaff}
              func={this.endActivation}
              force={this.forceSwitch}
            />
          );
        } else {
          machineBody = (
            <OccupiedMachine
              activation={this.props.activation}
              machine={this.props.machine}
              isStaff={isStaff}
              func={this.endActivation}
              force={this.forceSwitch}
            />
          );
        }
      } else {
        machineBody = (
          <FreeMachine
            machine={this.props.machine}
            isStaff={isStaff}
            func={this.startActivation}
            force={this.forceSwitch}
          />
        );
      }
    }
    var price;
    if (this.props.machine.Name.indexOf('Tutor') < 0) {
      price = ' [€';
      price += this.props.machine.Price.toFixed(2);
      price += '/';
      switch (this.props.machine.PriceUnit) {
        case 'hour':
          price += 'h';
          break;
        case 'minute':
          price += 'min';
          break;
        default:
          price += this.props.machine.PriceUnit;
      }
      price += ']';
    }
    var availabilityDisplay;
    if (_.isNumber(this.props.machine.ReservationPriceHourly)) {
      availabilityDisplay = <AvailabilityDisplay machineId={this.props.machine.Id}/>;
    }
    return (
      <div className="machine-container">
        <div className="container-fluid">
          <div className="machine-header">
            <div className="machine-title pull-left">{this.props.machine.Name} {price}</div>
            <div className="clearfix"></div>
          </div>
          <div className="machine-body">
            {machineBody}
            {availabilityDisplay}

            <ul className="machine-extra-actions">
              <li className="action-item">
                <MaintenanceSwitch machineId={this.props.machine.Id}/>
              </li>
              <li className="action-item">
                <RepairButton machineId={this.props.machine.Id}/>
              </li>
            </ul>

          </div>          
        </div>
      </div>
    );
  }
});

export default MachineChooser;
