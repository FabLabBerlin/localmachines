var ApiActions = require('./ApiActions');
var toastr = require('toastr');


var TutoringActions = {
	startTutoring(id) {
    ApiActions.postCall('/api/tutorings/' + id + '/start', {}, function() {
      toastr.info('Tutoring started');
    }, 'Cannot stop the tutoring');
  },

  stopTutoring(id) {
    ApiActions.postCall('/api/tutorings/' + id + '/stop', {}, function() {
      toastr.info('Tutoring stopped');
    }, 'Cannot start the tutoring');
  }
};

export default TutoringActions;
