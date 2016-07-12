var actions = require('../../modules/ForgotPassword/actions');
var { Navigation } = require('react-router');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


var Start = React.createClass({
  mixins: [ Navigation ],

  cancel() {
    this.replaceWith('/login');
  },

  handleSubmit(event) {
    event.preventDefault();
    var email = this.refs.email.value;
    if (email) {
      actions.emailReset(this.context.router, email);
    } else {
      toastr.error('Please enter an E-Mail address');
    }
    return false;
  },

  render() {
    return (
      <div className="container">
        <form onSubmit={this.handleSubmit}>
          <h3>What is your E-Mail address?</h3>
          <input
            ref="email"
            type="text"
            name="email"
            className="form-control"
            placeholder="E-Mail address"
            required={true}
            autoFocus="on"
            autoCorrect="off"
            autoCapitalize="off"
          />
          <hr/>
          <div className="pull-right">
            <button className="btn btn-info btn-lg wizard-button"
                onClick={this.cancel}>Cancel</button>
            <button className="btn btn-primary btn-lg wizard-button"
              type="submit">Reset password</button>
          </div>
        </form>
      </div>
    );
  }
});

export default Start;
