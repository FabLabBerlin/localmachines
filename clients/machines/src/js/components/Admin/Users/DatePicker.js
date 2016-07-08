var $ = require('jquery');
var ReactDatePicker = require('react-datepicker');
var moment = require('moment');
var React = require('react');


var DatePicker = React.createClass({
  componentDidMount() {
    $(this.refs.inputField.getDOMNode()).pickadate({
      format: 'yyyy-mm-dd'
    });
  },

  render() {
    return (
      <div className="input-group">
        <input ref="inputField"
               className="adm-user-membership-end-date form-control datepicker"
               placeholder={this.props.placeholder}
               defaultValue={moment(this.props.date).format('YYYY-MM-DD')}/>
        <div className="input-group-addon">
          <i className="fa fa-calendar"/>
        </div>
      </div>
    );
  }
});

export default DatePicker;
