var React = require('react');
var reactor = require('../../reactor');
var DTRDatePicker = require('../DateRangePicker/DTRDatePicker');

var ReservationsPage = React.createClass({

  getInitialState: function() {
    return {showDatePicker: false};
  },

  openDatePicker() {
    this.setState({showDatePicker: true});
  },

  closeDatePicker() {
    this.setState({showDatePicker: false});
  },

  render() {
    return (
      <div className="page-reservations">
        <div className="container">
          <h2>Reservations Page</h2>
          
          <button 
            className="btn btn-primary"
            onClick={this.openDatePicker}>
            Open Date Picker
          </button>
  
          {this.state.showDatePicker ? <DTRDatePicker 
            closeFunc={this.closeDatePicker} /> : ''}
            
          </div>
      </div>
    );
  }
});

export default ReservationsPage;
