var getters = require('../../../../getters');
var LoaderLocal = require('../../../LoaderLocal');
var Location = require('../../../../modules/Location');
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
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);
  },

  description() {
    const m = this.machine();

    if (!m) {
      return undefined;
    }

    return this.textToHTML(m.get('Description') || '');
  },

  links() {
    const m = this.machine();

    if (!m) {
      return undefined;
    }

    var ls = m.get('Links') || '';
    return (
      <div>
        {ls.split('\n').map((line, i) => {
          if (!line.trim()) {
            return undefined;
          }

          const tmp = line.split(' ');
          const href = tmp[0].startsWith('http')
            ? tmp[0]
            : 'http://' + tmp[0].trim();
          const label = tmp.length === 1
            ? href
            : tmp.slice(1).join(' ');

          return (
            <p key={i}>
              <a href={href}>{label}</a>
            </p>
          );
        })}
      </div>
    );
  },

  machine() {
    const machineId = parseInt(this.props.params.machineId);
    var m;

    if (this.state.machines) {
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === machineId;
      });
    }

    return m;
  },

  render() {
    const m = this.machine();

    if (!m) {
      return <LoaderLocal/>;
    }

    return (
      <div id="m-info" className="row">
        <div id="m-info-left" className="col-md-6 col-md-push-6">
          <Section id="m-info-safety-guidelines" title="Safety Guidelines">
            <table>
              <tbody>
                <tr>
                  <td>{this.safetyGuidelines()}</td>
                </tr>
              </tbody>
            </table>
          </Section>

          <Section id="m-info-links" title="Links">
            <table>
              <tbody>
                <tr>
                  <td>{this.links()}</td>
                </tr>
              </tbody>
            </table>
          </Section>
        </div>

        <div id="m-info-right" className="col-md-6 col-md-pull-6">
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

          <Section id="m-info-description" title="Description">
            <table>
              <tbody>
                <tr>
                  <td>{this.description()}</td>
                </tr>
              </tbody>
            </table>
          </Section>
        </div>
      </div>
    );
  },

  safetyGuidelines() {
    const m = this.machine();

    if (!m) {
      return undefined;
    }

    return this.textToHTML(m.get('SafetyGuidelines') || '');
  },

  textToHTML(text) {
    return (
      <div>
        {text.split('\n').map((line, i) => {
          return <p key={i}>{line}</p>;
        })}
      </div>
    );
  }
});

export default InfoPage;
