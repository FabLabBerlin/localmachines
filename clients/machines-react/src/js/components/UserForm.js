import React from 'react';
import _ from 'lodash';

/*
 * Import toastr and set position
 */
import toastr from 'toastr';
toastr.options.possitionClass = 'toast-bottom-left';

/*
 * UserForm component:
 * manage the form the user fill to update his profile
 */
var UserForm = React.createClass({

  /*
   * Handle the change in the input
   * @event: typing in the form
   * User function passed by props
   */
  handleChangeForm(event) {
    event.preventDefault();
    this.props.func(event);
  },

  /*
   * Ask the UserPage the send the form to the server
   * Use function passed by props
   */
  handleSubmit(event) {
    event.preventDefault();
    this.props.submit();
  },

  /*
   * Update the user password
   * make syntax verification
   * if all good, give the upper hand to the UserPage
   */
  updatePassword(event) {
    event.preventDefault();
    var minPassLength = 3;
    var password = document.getElementById('password');
    if(password.value !== document.getElementById('repeat').value){
      toastr.error('Passwords do not match');
    } else if(!password.value || password.value === '') {
      toastr.error('You did not write any new password');
    } else if(password.value.length < minPassLength) {
      toastr.error('Password too short');
    } else {
      this.props.passwordFunc(password.value);
    }
  },

  /*
   * Render the form:
   * for each information in the userPage state:
   *  - create an input
   *  - give the value of the state to the related input
   */
  render() {
    var NodeInput = _.map(this.props.info, function(value, key) {
      return (
        <div className="col-md-6" key={key}>
          <div className="form-group">
            <label htmlFor="user-information">{key}</label>
            <input type="text" value={value} 
              id={key}
              className="form-control"
              onChange={this.handleChangeForm}
            />
          </div>
        </div>
      );
    }.bind(this));
        return (
            <div className="userForm" >
                <form onSubmit={this.handleSubmit} >
                    {NodeInput}
                    <br />
                    <input type="password"
                    />
                    <input type="password"
                    />
                    <button>Okay</button>
                </form>
            </div>
        );
    }
});

module.exports = UserForm;
