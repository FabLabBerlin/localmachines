(function(){

'use strict';

var mod = angular.module("fabsmith.admin.randomtoken", []);

mod.service('randomToken', function() {
  var tokens = [
    'Randy3time',
    'Token2be4me',
    'Token4life',
    'Randomi7er',
    'RandomSk8ter',
    'H8tersGonn4'
  ];
  var id = Math.round(Math.random() * (tokens.length-1));

  this.generate = function() {
    return tokens[id];
  };

  return this;
});

})(); // closure