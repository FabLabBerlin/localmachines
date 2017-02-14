import $ from 'jquery';
import Categories from '../../../modules/Categories';
import getters from '../../../getters';
import LoaderLocal from '../../LoaderLocal';
import Location from '../../../modules/Location';
import React from 'react';
import reactor from '../../../reactor';
import TableCRUD from '../../TableCRUD/TableCRUD';
import toastr from 'toastr';
import UserActions from '../../../actions/UserActions';


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

  add() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    $.ajax({
      url: '/api/machine_types?location=' + locationId,
      dataType: 'json',
      type: 'POST',
      data: {
        name: 'Untitled'
      }
    })
    .done(() => {
      toastr.info('Successfully added category.');
      this.fetchData();
    })
    .fail(() => {
      toastr.error('Error adding category.  Please try again later.');
    });
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
                      onAdd={this.add}
                      onAfterUpdate={this.fetchData}
                      updateUrl="/api/machine_types"/>;
  }

});

export default CategoriesPage;
