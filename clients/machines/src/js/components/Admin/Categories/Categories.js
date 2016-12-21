var Categories = require('../../../modules/Categories');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var Location = require('../../../modules/Location');
var React = require('react');
var reactor = require('../../../reactor');
var TableCRUD = require('../../TableCRUD/TableCRUD');
var UserActions = require('../../../actions/UserActions');


var CategoriesPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      categories: Categories.getters.getAll,
      locationId: Location.getters.getLocationId
    };
  },

  componentWillMount() {
    this.fetchData();
  },

  fetchData() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);

    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);
    Categories.actions.loadAll(locationId);
  },

  render() {
    if (!this.state.categories) {
      return <LoaderLocal/>;
    }

    const fields = [
      {key: 'Id', label: 'Id'},
      {key: 'ShortName', label: 'Short Name'},
      {key: 'Name', label: 'Name'}
    ];

    return <TableCRUD entities={this.state.categories}
                      fields={fields}
                      onAfterUpdate={this.fetchData}
                      updateUrl="/api/machine_types"/>;
  }

});

export default CategoriesPage;
