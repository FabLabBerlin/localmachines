var $ = require('jquery');
import actionTypes from './actionTypes';
import React from 'react';
import reactor from '../../reactor';
import toastr from '../../toastr';


var CategoriesActions = {
  loadAll(locId) {
    $.ajax({
      url: '/api/machine_types',
      data: {
        location: locId
      }
    })
    .success(categories => {
      reactor.dispatch(actionTypes.SET_CATEGORIES, categories);
    })
    .error(() => {
      toastr.error('Error fetching categories.');
    });
  }
};

export default CategoriesActions;
