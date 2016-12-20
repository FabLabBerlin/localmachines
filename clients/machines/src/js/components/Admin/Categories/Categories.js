var Categories = require('../../../modules/Categories');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var Location = require('../../../modules/Location');
var React = require('react');
var reactor = require('../../../reactor');
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

    return (
      <div className="container">
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Category Id</th>
              <th>Short Name</th>
              <th>Name</th>
            </tr>
          </thead>

          <tbody>
            {this.state.categories.map((c, i) => {
              return (
                <tr key={i}>
                  <td>{c.get('Id')}</td>
                  <td>{c.get('ShortName')}</td>
                  <td>{c.get('Name')}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  }

});

export default CategoriesPage;
