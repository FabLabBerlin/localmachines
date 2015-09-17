var FeedbackDialogs = require('../Feedback/FeedbackDialogs');
var React = require('react');
var LoginActions = require('../../actions/LoginActions');
var MachineActions = require('../../actions/MachineActions');
var OccupiedMachine = require('./OccupiedMachine');
var BusyMachine = require('./BusyMachine');
var FreeMachine = require('./FreeMachine');

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

  /*
   * To force a switch
   * Only admin have to be able to use this function
   */
  forceSwitch(onOrOff) {
    let mid = this.props.info.Id;
    let aid = this.props.activation.Id;
    if(onOrOff === 'off') {
      MachineActions.adminTurnOffMachine(mid, aid);
    } else if (onOrOff === 'on') {
      MachineActions.adminTurnOnMachine(mid);
    }
  },

  /*
   * Function pass by props to children
   * End an activation
   */
  endActivation() {
    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to stop the activation for <b>' +
        this.props.info.Name + '</b>?',
      callback: function(confirmed) {
        if (confirmed) {
          let aid = this.props.activation.Id;
          MachineActions.endActivation(aid, function() {
            //FeedbackDialogs.checkSatisfaction(aid);
          }.bind(this));
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }.bind(this)
    });

    LoginActions.keepAlive();
  },

  /*
   * Function pass by props to children
   * Start an activation
   */
  startActivation() {
    let mid = this.props.info.Id;
    MachineActions.startActivation(mid);
    LoginActions.keepAlive();
  },

  /*
   * If the id of the activation is the same, doesn't render the component
   * To not rerender each component at each polling
   */
  shouldComponentUpdate(nextProps) {
    var shouldUpdate = true;
    if( nextProps.activation.Id === this.props.activation.Id && !this.props.activation.FirstName ) {
      shouldUpdate = false;
    }
    return shouldUpdate;
  },

  /*
   * Render component
   * Can choose what component will be display depending on the props
   * @busy + @sameUser => BusyMachine
   * @busy => OccupiedMachine
   * @nothing => FreeMachine
   */
  render() {
    let isAdmin = this.props.user.get('UserRole') === 'admin';
    return (
      <div className="machine available">
        <div className="container-fluid">
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
                  isAdmin={isAdmin}
                  func={this.endActivation}
                  force={this.forceSwitch}
                />
            ) : (
            <OccupiedMachine
              activation={this.props.activation}
              info={this.props.info}
              isAdmin={isAdmin}
              func={this.endActivation}
              force={this.forceSwitch}
            />
            ) :
              (
                <FreeMachine
                  info={this.props.info}
                  isAdmin={isAdmin}
                  func={this.startActivation}
                  force={this.forceSwitch}
                />
            )}
          </div>
        </div>
      </div>
    );
  }
});

export default MachineChooser;
