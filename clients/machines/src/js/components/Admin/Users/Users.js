import _ from 'lodash';
var $ = require('jquery');
import getters from '../../../getters';
import LoaderLocal from '../../LoaderLocal';
import LocationGetters from '../../../modules/Location/getters';
import MachineActions from '../../../actions/MachineActions';
import { Navigation } from 'react-router';
import React from 'react';
import reactor from '../../../reactor';
import Users from '../../../modules/Users';


var User = React.createClass({
  render() {
    const user = this.props.user;

    return (
      <tr>
        <td>{user.FirstName}</td>
        <td>{user.LastName}</td>
        <td>{user.Username}</td>
        <td>{user.Email}</td>
        <td>
          {user.NoAutoInvoicing ?
            <i className="fa fa-check"/> : null
          }
        </td>
        <td>{user.Company}</td>
        <td>
          <a className="btn btn-primary btn-ico pull-right" 
          href={'/machines/#/admin/users/' + user.Id}>
            <span className="fa fa-edit"></span>
          </a>
        </td>
      </tr>
    );
  }
});


var UsersView = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      users: Users.getters.getUsers
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    Users.actions.fetchUsers({locationId});
  },

  render() {
    if (this.state.users) {
      var users = _.sortBy(this.state.users, m => m.Name);
      console.log('users=', users);
      return (
        <div className="container-fluid">
          <table className="table table-striped table-hover">
            <thead>
              <tr>
                <th>First Name</th>
                <th>Last Name</th>
                <th>Username</th>
                <th>Email</th>
                <th>No Auto Invoicing</th>
                <th>Company</th>
                <th>&nbsp;</th>
              </tr>
            </thead>
            <tbody>
              {_.map(users, (u, i) => {
                return <User key={i} user={u}/>;
              })}
            </tbody>
          </table>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }

});

export default UsersView;
