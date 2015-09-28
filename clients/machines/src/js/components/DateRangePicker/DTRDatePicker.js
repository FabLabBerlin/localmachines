var React = require('react');

var DTRDatePicker = React.createClass({
  
  getInitialState: function() {
    return {dayOfWeek: 0};
  },

  componentDidMount() {
    // Create an array with months
    var today = new Date();

    // We better work with UTC dates here 
    var dayOfWeek = today.getDay();

    var date = today.getDate();

    // Populate months starting from today
    /*
    var displayNumMonths = 6;
    var currentMonth = today.getMonth();
    for (var i=0; i<displayNumMonths; i++) {
      var loopDate = new Date()
    }
    */

    this.setState({dayOfWeek: dayOfWeek});
  },

  onCancel() {
    this.props.closeFunc();
  },

  render() {

    return (
      <div className="range-date-picker">

        <div className="header">
          <div className="container">
            <h2>Select date</h2>
            <p>Day of week: {this.state.dayOfWeek}</p>
          </div>
        </div>  

        <div className="content">
          <div className="container">
            <table className="table table-striped">
              <caption>September 2015</caption>
              <thead>
                <tr>
                  <th>Mo</th>
                  <th>Tu</th>
                  <th>We</th>
                  <th>Th</th>
                  <th>Fr</th>
                  <th>Sa</th>
                  <th>Su</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>1</td>
                  <td>2</td>
                  <td>3</td>
                  <td>4</td>
                  <td>5</td>
                  <td>6</td>
                  <td>7</td>
                </tr>

              </tbody>
            </table>
          </div>
        </div>

        <div className="footer">
          <div className="container">
            <button 
              className="btn btn-danger"
              onClick={this.onCancel}>
              Cancel
            </button>

            <button 
              className="btn btn-primary"
              onClick={this.onCancel}>
              OK
            </button>
          </div>
        </div>

      </div>
    );
  }
});

export default DTRDatePicker;