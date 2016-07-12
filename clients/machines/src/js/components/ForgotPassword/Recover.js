var actions = require('../../modules/ForgotPassword/actions');
var { Navigation } = require('react-router');
var React = require('react');
var toastr = require('../../toastr');


var Recover = React.createClass({
  mixins: [ Navigation ],

  handleSubmit(event) {
    event.preventDefault();
    var phone = this.refs.phone.value;
    if (phone) {
      actions.submitPhone(this.context.router, phone);
    } else {
      toastr.error('Please enter your phone number');
    }
    return false;
  },

  render() {
    return (
      <div className="container">
        <form onSubmit={this.handleSubmit}>
          <h3>What is your phone number?</h3>
          <input
            ref="phone"
            type="text"
            name="phone"
            className="form-control"
            placeholder="Phone number"
            required={true}
            autofocus
            autoCorrect="off"
            autoCapitalize="off"
          />
          <hr/>
          <div className="pull-right">
            <button className="btn btn-primary btn-lg wizard-button"
              type="submit">Submit</button>
          </div>
        </form>
      </div>
    );
  }
});

export default Recover;
