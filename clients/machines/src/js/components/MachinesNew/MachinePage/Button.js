var $ = require('jquery');
var ActivationTimer = require('../../MachinePage/Machine/ActivationTimer');
var constants = require('../constants');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LoginActions = require('../../../actions/LoginActions');
var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../reactor');
var ReservationTimer = require('../../MachinePage/Machine/ReservationTimer');


// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';


var Button = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      uid: getters.getUid
    };
  },

  activationEnd() {
    let aid = this.props.machine.getIn(['activation', 'Id']);

    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to stop the activation for <b>' +
        this.props.machine.get('Name') + '</b>?',
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

  activationStart() {
    const mid = this.props.machine.get('Id');
    MachineActions.startActivation(mid);
    LoginActions.keepAlive();
  },

  render() {
    switch (this.props.status) {
    case constants.AVAILABLE:
      return (
        <div className="m-action"
             onClick={this.activationStart}>
          START
        </div>
      );
    case constants.LOCKED:
      return (
        <div className="m-action">
          LOCKED
        </div>
      );
    case constants.MAINTENANCE:
      return (
        <div className="m-action">
          {this.props.isStaff ? (
            this.props.machine.get('activation') ? (
              <div onClick={this.activationEnd} style={{cursor: 'pointer'}}>
                STOP
              </div>
            ) : (
              <div onClick={this.activationStart} style={{cursor: 'pointer'}}>
                START
              </div>
            )
          ) : (
            <span>MAINTENANCE</span>
          )}
        </div>
      );
    case constants.OCCUPIED:
      if (this.props.isStaff) {
        return (
          <div className="m-action m-clock"
               onClick={this.activationEnd}>
            <span>STOP</span>
            <ActivationTimer activation={this.props.machine.get('activation').toJS()}/>
          </div>
        );
      }
      return (
        <div className="m-action">
          <span>OCCUPIED</span>
          <ActivationTimer activation={this.props.machine.get('activation').toJS()}/>
        </div>
      );
    case constants.RESERVED:
      return (
        <div className="m-action m-clock">
          {this.props.isStaff ||
           this.props.reservation.get('UserId') === this.state.uid
          ? (
            <div onClick={this.activationStart} style={{cursor: 'pointer'}}>
              START
            </div>
          ) : (
            <div>
              RESERVED
              <ReservationTimer reservation={this.props.reservation.toJS()}/>
            </div>
          )}
        </div>
      );
    case constants.RUNNING:
      return (
        <div className="m-action m-clock"
             onClick={this.activationEnd}>
          STOP
          <ActivationTimer activation={this.props.machine.get('activation').toJS()}/>
        </div>
      );
    default:
      console.log('this.props.status=', this.props.status);
      return <LoaderLocal/>;
    }
  }
});

export default Button;
