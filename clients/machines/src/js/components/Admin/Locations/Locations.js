var $ = require('jquery');
import Button from '../../Button';
import LoaderLocal from '../../LoaderLocal';
import Location from '../../../modules/Location';
import React from 'react';
import getters from '../../../getters';
import reactor from '../../../reactor';
import TableCRUD from '../../TableCRUD/TableCRUD';
import toastr from 'toastr';
import UserActions from '../../../actions/UserActions';


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
    this.fetchData();
  },

  add() {
    $.ajax({
      url: '/api/locations',
      dataType: 'json',
      type: 'POST',
      data: {
        title: 'Untitled'
      }
    })
    .done(() => {
      toastr.info('Successfully added location.');
      Location.actions.loadLocations();
    })
    .fail(() => {
      toastr.error('Error adding location.  Please try again later.');
    });
  },

  fetchData() {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    Location.actions.loadLocations();
    Location.actions.loadUserLocations(uid);
  },

  render() {
    if (!this.state.user || !this.state.locations) {
      return <LoaderLocal/>;
    }

    if (!this.state.user.get('SuperAdmin')) {
      return <div/>;
    }

    const fields = [
      {key: 'Id', label: 'Id'},
      {key: 'Title', label: 'Title'},
      {key: 'FirstName', label: 'First Name'},
      {key: 'LastName', label: 'Last Name'},
      {key: 'Email', label: 'E-Mail'},
      {key: 'XmppId', label: 'Jabber ID'},
      {key: 'City', label: 'City'},
      {key: 'Timezone', label: 'IANA Timezone'},
      {key: 'Approved', label: 'Show on Login page'}
    ];

    return <TableCRUD entities={this.state.locations}
                      fields={fields}
                      onAdd={this.add}
                      onAfterUpdate={this.fetchData}
                      updateUrl="/api/locations"/>;
  }
});

export default Locations;
