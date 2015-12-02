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
      tutorings = _.map(this.state.tutorings.toJS(), function(t, i) {
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
