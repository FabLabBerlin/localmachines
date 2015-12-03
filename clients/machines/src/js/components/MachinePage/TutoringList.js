var getters = require('../../getters');
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
      tutorings = _.sortBy(this.state.tutorings.toJS(), function(t) {
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
  }
});

export default TutoringList;
