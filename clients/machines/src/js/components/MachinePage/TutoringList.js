var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var Tutoring = require('./Tutoring');

var TutoringList = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      tutorings: getters.getTutorings
    };
  },

  render() {
    var tutorings = [];

    if (this.state.tutorings) {
      const today = moment().format('YYYY-MM-DD');
      tutorings = _.filter(this.state.tutorings.toJS(), function(t) {
        if (t.TimeStart) {
          return moment(t.TimeStart).format('YYYY-MM-DD') === today;
        }
      });
      tutorings = _.sortBy(tutorings, function(t) {
        if (t.Running) {
          return -1;
        } else {
          return 0;
        }
      });
      tutorings = _.map(tutorings, function(t, i) {
        return <Tutoring key={i}
                         tutoring={t}/>;
      });
    }

    if (tutorings.length > 0) {
      return (
        <div className="tutoring-list">
          
          <div className="tutoring-list-header">
            <div className="container-fluid">
              <h2>Your tutorings today</h2>
            </div>
          </div>
          
          <div className="tutoring-list-body">
            {tutorings}
          </div>
        
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default TutoringList;
