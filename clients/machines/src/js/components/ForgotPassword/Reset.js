import actions from '../../modules/ForgotPassword/actions';
import React from 'react';
import toastr from '../../toastr';


var Reset = React.createClass({
  handleSubmit(event) {
    event.preventDefault();
    var pass = this.refs.password.value;
    var repeat = this.refs.repeat.value;
    if (pass && repeat) {
      if (pass !== repeat) {
        toastr.error('Both passwords must match');
        return;
      }
      if (pass.length < 6) {
        toastr.error('Password length must be at least 6');
        return;
      }
      actions.submitPassword(this.context.router, pass);
    } else {
      toastr.error('Please enter your new password and repeat it');
    }
  },

  render() {
    return (
      <div className="container">
        <form onSubmit={this.handleSubmit}>
          <h3>Please enter a new password</h3>
          <input
            ref="password"
            type="password"
            name="password"
            className="form-control"
            placeholder="Password"
            required={true}
            autofocus
            autoCorrect="off"
            autoCapitalize="off"
          />
          <input
            ref="repeat"
            type="password"
            name="repeat"
            className="form-control"
            placeholder="Repeat"
            required={true}
            autofocus
            autoCorrect="off"
            autoCapitalize="off"
          />
          <hr/>
          <div className="pull-right">
            <button className="btn btn-primary btn-lg wizard-button"
              type="submit">Reset password</button>
          </div>
        </form>
      </div>
    );
  }
});

export default Reset;
