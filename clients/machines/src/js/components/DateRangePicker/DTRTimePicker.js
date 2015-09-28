var React = require('react');

var DTRDatePicker = React.createClass({
  
  onCancel() {
    this.setState({showDatePicker: false});
  },

  render() {
    return (
      <div className="range-date-picker">
        Date/time range time picker

        <button 
          className="btn btn-danger"
          onClick={this.onCancel}>
          Cancel
        </button>
      </div>
    );
  }
});