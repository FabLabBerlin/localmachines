var Button = require('../../Button');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationActions = require('../../../actions/LocationActions');
var LocationGetters = require('../../../modules/Location/getters');
var Memberships = require('../../../modules/Memberships');
var React = require('react');
var reactor = require('../../../reactor');
var SettingsActions = require('../../../modules/Settings/actions');
var UserActions = require('../../../actions/UserActions');


var MembershipsPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      memberships: Memberships.getters.getAllMemberships
    };
  },

  componentDidMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    Memberships.actions.fetch({locationId});
    LocationActions.loadUserLocations(uid);
    SettingsActions.loadSettings({locationId});
  },

	render() {
    if (!this.state.memberships) {
      return <LoaderLocal/>;
    }

    return (
      <div className="container-fluid">
        <div className="row">
          <div className="col-xs-6 text-left">
            <h1>All Memberships</h1>
          </div>
          <div className="col-xs-6 text-right">
            <Button.Annotated id="mbs-toggle-inactive"
                               label="Show archived"
                               onClick={this.setShowArchived.bind(this, true)}/>
          </div>
        </div>
        <table>
          <thead>
            <tr>
              <th>
                Name
              </th>
              <th>
                Minimum Duration
              </th>
              <th>
                Monthly Price
              </th>
            </tr>
          </thead>
          <tbody>
            {this.state.memberships.sortBy(mb => mb.get('Title'))
                                   .map(mb => {
              return (
                <tr key={mb.get('Id')}>
                  <td>{mb.get('Title')}</td>
                  <td>{mb.get('DurationMonths')}</td>
                  <td>{mb.get('MonthlyPrice')}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  },

  setShowArchived(yes, e) {
  }
});

export default MembershipsPage;
