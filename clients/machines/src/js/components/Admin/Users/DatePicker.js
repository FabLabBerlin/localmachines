var React = require('react');


var DatePicker = React.createClass({
  render() {
    return (
      <div className="input-group">
        <input
          className="adm-user-membership-end-date form-control datepicker"
          placeholder={this.props.placeholder}
          value={this.props.date}/>
        <div className="input-group-addon">
          <i className="fa fa-calendar"/>
        </div>
      </div>
    );
  }
});

export default DatePicker;
