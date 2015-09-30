var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');


var DatePicker = React.createClass({
  render() {
    return (
      <div>
        <h3 className="h3">Select Date</h3>
        <input type="text" placeholder="YYYY-MM-DD" ref="date"/>
        <button className="btn btn-lg btn-primary" type="button" onClick={this.setDate}>Next</button>
      </div>
    );
  },

  setDate() {
    var date = this.refs.date.getDOMNode().value;
    ReservationsActions.createSetDate({ date });
  }
});

export default DatePicker;
