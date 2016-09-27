var $ = require('jquery');
var constants = require('../constants');
var FeedbackDialogs = require('../../Feedback/FeedbackDialogs');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationActions = require('../../../actions/LocationActions');
var LocationGetters = require('../../../modules/Location/getters');
var LoginActions = require('../../../actions/LoginActions');
var MachineActions = require('../../../actions/MachineActions');
var MachineMixin = require('../MachineMixin');
var Machines = require('../../../modules/Machines');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('../../../toastr');
var UserActions = require('../../../actions/UserActions');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';


var Button = React.createClass({
  /*
   * Function pass by props to children
   * End an activation
   */
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

  /*
   * Function pass by props to children
   * Start an activation
   */
  activationStart() {
    const mid = this.props.machine.get('Id');
    MachineActions.startActivation(mid);
    LoginActions.keepAlive();
  },

  render() {
    switch (this.props.status) {
    case constants.AVAILABLE:
      return (
        <div className="m-action m-start"
             onClick={this.activationStart}>
          START
        </div>
      );
    case constants.RUNNING:
      return (
        <div className="m-action m-stop"
             onClick={this.activationEnd}>
          STOP
        </div>
      );
    default:
      console.log('this.props.status=', this.props.status);
      return <LoaderLocal/>;
    }
  }
});


var MachinePage = React.createClass({

  mixins: [ MachineMixin, reactor.ReactMixin ],

  getDataBindings() {
    return {
      activations: Machines.getters.getActivations,
      locationId: LocationGetters.getLocationId,
      machines: Machines.getters.getMachines,
      user: getters.getUser
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    LocationActions.loadUserLocations(uid);
    MachineActions.wsDashboard(null, locationId);
  },

  machine() {
    const machineId = parseInt(this.props.params.machineId);
    var m;

    if (this.state.machines) {
      console.log('this.state.machines=', this.state.machines);
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === machineId;
      });

      if (this.state.activations) {
        console.log('this.state.activations->true');
        const as = this.state.activations
        .groupBy(a => a.get('MachineId'))
        .get(m.get('Id'));

        if (as) {
          console.log('as->true');
          m = m.set('activation', as.get(0));
        }
      }
    }

    return m;
  },

  render() {
    const m = this.machine();
    var button;

    if (!m) {
      return <LoaderLocal/>;
    }

    switch (this.status()) {
      case constants.AVAILABLE:
        button = (
          <div className="m-action m-start">
            START
          </div>
        );
        break;
      case constants.RUNNING:
        button = (
          <div className="m-action m-stop">
            STOP
          </div>
        );
        break;
    }

    const small = window.innerWidth < 500;
    const style = {
      backgroundImage: 'url(' + this.imgUrl(small) + ')'
    };

    return (
      <div className="container-fluid">
        <div id="m-header">
          <h2>{m.get('Name')} ({m.get('Brand')})</h2>
          <div id="m-img" style={style}/>
          <div id="m-header-panel">
            <Button machine={this.machine()}
                    status={this.status()}/>
          </div>
          <div id="m-report" onClick={this.repair}>
            <span>Report a machine failure</span>
          </div>
        </div>
      </div>
    );
  },

  repair() {
    FeedbackDialogs.machineIssue(this.machine().get('Id'));
  }

});

export default MachinePage;
