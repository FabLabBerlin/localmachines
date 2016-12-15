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

  handleEdit(key, e) {
    Location.actions.setEditLocation({[key]: e.target.value});
  },

  handleSave() {
    Location.actions.saveEditedLocation();
  },

  handleSelect() {
    Location.actions.setEditLocation(this.props.location.toJS());
  },

  render() {
    const l = this.props.location;
    const edit = this.state.editLocation;

    return (
      edit.get('Id') !== l.get('Id') ? (
        <tr onClick={this.handleSelect}>
          <td>{l.get('Id')}</td>
          <td>{l.get('Title')}</td>
          <td>{l.get('FirstName')}</td>
          <td>{l.get('LastName')}</td>
          <td>{l.get('Email')}</td>
          <td>{l.get('XmppId')}</td>
          <td>{l.get('City')}</td>
          <td>{l.get('Timezone')}</td>
          <td/>
        </tr>
      ) : (
        <tr>
          <td>{l.get('Id')}</td>
          <td><input onChange={this.handleEdit.bind(this, 'Title')}
                     value={edit.get('Title')}/></td>
          <td><input onChange={this.handleEdit.bind(this, 'FirstName')}
                     value={edit.get('FirstName')}/></td>
          <td><input onChange={this.handleEdit.bind(this, 'LastName')}
                     value={edit.get('LastName')}/></td>
          <td><input onChange={this.handleEdit.bind(this, 'Email')}
                     value={edit.get('Email')}/></td>
          <td><input onChange={this.handleEdit.bind(this, 'XmppId')}
                     value={edit.get('XmppId')}/></td>
          <td><input onChange={this.handleEdit.bind(this, 'City')}
                     value={edit.get('City')}/></td>
          <td><input onChange={this.handleEdit.bind(this, 'Timezone')}
                     value={edit.get('Timezone')}/></td>
          <td><i className="fa fa-floppy-o"
                 onClick={this.handleSave}/></td>
        </tr>
      )
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
              <th/>
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
