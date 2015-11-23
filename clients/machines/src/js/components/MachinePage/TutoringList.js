var React = require('react');
var Tutoring = require('./Tutoring');

var TutoringList = React.createClass({
  render() {
    var tutorings = [];
    tutorings.push(<Tutoring />);
    tutorings.push(<Tutoring />);

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
