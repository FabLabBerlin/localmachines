var LoaderLocal = require('../../LoaderLocal');
var Location = require('../../../modules/Location');
var React = require('react');
var getters = require('../../../getters');
var reactor = require('../../../reactor');
var UserActions = require('../../../actions/UserActions');


var Row = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      editLocation: Location.getters.getEditLocation
    };
  },

  handleSelect() {
    const id = this.props.location.get('Id');
    Location.actions.setEditLocation({Id: id});
  },

  render() {
    const l = this.props.location;

    return (
      <tr onClick={this.handleSelect}>
        <td>{l.get('Id')}</td>
        <td>{l.get('Title')}</td>
        <td>{l.get('FirstName')}</td>
        <td>{l.get('LastName')}</td>
        <td>{l.get('Email')}</td>
        <td>{l.get('XmppId')}</td>
        <td>{l.get('City')}</td>
        <td>{l.get('Timezone')}</td>
      </tr>
    );
  }
});


var Locations = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      editLocation: Location.getters.getEditLocation,
      isLogged: getters.getIsLogged,
      locations: Location.getters.getLocations,
      user: getters.getUser
    };
  },

  componentDidMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);
  },

  render() {
    if (!this.state.user || !this.state.locations) {
      return <LoaderLocal/>;
    }

    if (!this.state.user.get('SuperAdmin')) {
      return <div/>;
    }

    console.log('locations=', this.state.locations);
    console.log('editLocation=', this.state.editLocation.toJS());

    return (
      <div className="container">
        <h1>Location List</h1>
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Id</th>
              <th>Title</th>
              <th>First Name</th>
              <th>Last Name</th>
              <th>E-Mail</th>
              <th>Jabber ID</th>
              <th>City</th>
              <th>IANA Timezone</th>
            </tr>
          </thead>
          <tbody>
            {this.state.locations.map((l, i) => <Row key={i} location={l}/>)}
          </tbody>
        </table>
      </div>
    );
  }
});

export default Locations;
