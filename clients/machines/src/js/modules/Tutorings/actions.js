var ApiActions = require('../Api/actions');
var getters = require('./getters');
var reactor = require('../../reactor');
var toastr = require('toastr');


var TutoringActions = {
	startTutoring(id) {
    ApiActions.postCall('/api/tutorings/' + id + '/start', {}, function() {
      toastr.info('Tutoring started');
    }, 'Cannot start the tutoring');
  },

  stopTutoring(id) {
    ApiActions.postCall('/api/tutorings/' + id + '/stop', {}, function() {
      toastr.info('Tutoring stopped');
    }, 'Cannot stop the tutoring');
  }
};

export default TutoringActions;
