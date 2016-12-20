var $ = require('jquery');
var actionTypes = require('./actionTypes');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


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
