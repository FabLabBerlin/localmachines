var getters = require('../../../../getters');
var LoaderLocal = require('../../../LoaderLocal');
var LocationActions = require('../../../../actions/LocationActions');
var LocationGetters = require('../../../../modules/Location/getters');
var MachineActions = require('../../../../actions/MachineActions');
var Machines = require('../../../../modules/Machines');
var React = require('react');
var reactor = require('../../../../reactor');
var UserActions = require('../../../../actions/UserActions');


var Section = React.createClass({
  render() {
    return (
      <div id={this.props.id} className="m-info-section">
        <div className="m-info-section-title">
          {this.props.title}
        </div>
        <hr/>
        <div className="m-info-section-content">
          {this.props.children}
        </div>
      </div>
    );
  }
});


var InfoPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    LocationActions.loadUserLocations(uid);
  },

  render() {
    const machineId = parseInt(this.props.params.machineId);
    var m;

    if (this.state.machines) {
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === machineId;
      });
    }

    if (!m) {
      return <LoaderLocal/>;
    }

    return (
      <div id="m-info" className="container-fluid">
        <Section id="m-info-specs" title="Technical Specifications">
          <table>
            <tbody>
              <tr>
                <td>Build volume:</td>
                <td>{m.get('WorkspaceDimensions')}</td>
              </tr>
            </tbody>
          </table>
        </Section>
      </div>
    );
  }
});

export default InfoPage;
