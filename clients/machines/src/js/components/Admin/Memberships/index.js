var Button = require('../../Button');
var getters = require('../../../getters');
var {hashHistory} = require('react-router');
var LoaderLocal = require('../../LoaderLocal');
var LocationActions = require('../../../actions/LocationActions');
var LocationGetters = require('../../../modules/Location/getters');
var Memberships = require('../../../modules/Memberships');
var React = require('react');
var reactor = require('../../../reactor');
var Settings = require('../../../modules/Settings');
var UserActions = require('../../../actions/UserActions');


var MembershipsPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      currency: Settings.getters.getCurrency,
      memberships: Memberships.getters.getAllMemberships,
      showArchived: Memberships.getters.getShowArchived
    };
  },

  componentDidMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    Memberships.actions.fetch({locationId});
    LocationActions.loadUserLocations(uid);
    Settings.actions.loadSettings({locationId});
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
            {this.state.showArchived ?
              <Button.Annotated id="mbs-toggle-inactive"
                                 label="Hide archived"
                                 onClick={this.setShowArchived.bind(this, false)}/> :
              <Button.Annotated id="mbs-toggle-inactive"
                                 label="Show archived"
                                 onClick={this.setShowArchived.bind(this, true)}/>
            }
          </div>
        </div>
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>
                Name
              </th>
              <th>
                Minimum Duration
              </th>
              <th>
                Monthly Price / {this.state.currency}
              </th>
            </tr>
          </thead>
          <tbody>
            {this.state.memberships.sortBy(mb => mb.get('Title'))
                                   .filter(mb => !mb.get('Archived')
                                                 || this.state.showArchived)
                                   .map(mb => {
              return (
                <tr key={mb.get('Id')} onClick={this.showMembership.bind(this, mb.get('Id'))}>
                  <td><b>{mb.get('Title')}</b></td>
                  <td>
                    {mb.get('DurationMonths')} {mb.get('DurationMonths') > 1 ?
                                                'months' : 'month'}
                  </td>
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
    Memberships.actions.setShowArchived(yes);
  },

  showMembership(id) {
    hashHistory.push('/admin/memberships/' + id);
  }
});

export default MembershipsPage;
