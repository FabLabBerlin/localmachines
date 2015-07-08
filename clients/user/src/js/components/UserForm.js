/*
 * Dependencies for toastr
 */
import ReactToastr from 'react-toastr';
var {ToastContainer} = ReactToastr;
var ToastMessageFactory = React.createFactory(ReactToastr.ToastMessage.animation);
import _ from 'lodash';
import React from 'react';

/*
 * UserForm component:
 * manage the form the user fill to update his profile
 */
var UserForm = React.createClass({

  addAlert(title, message) {
    this.clearAlert();
    this.refs.container.success(message, title, {
                                  closeButton: true,
                                  timeOut: 0,
                                  extendedTimeOut: 0
                                });
  },

  clearAlert() {
    this.refs.container.clear();
  },

  /*
   * Handle the change in the input
   * @event: typing in the form
   * User function passed by props
   */
  handleChangeForm(event) {
    this.props.func(event);
  },

  /*
   * Ask the UserPage the send the form to the server
   * Use function passed by props
   */
  handleSubmit() {
    this.props.submit();
  },

  /*
   * Update the user password
   * make syntax verification
   * if all good, give the upper hand to the UserPage
   */
  updatePassword() {
    var minPassLength = 3;
    var password = document.getElementById('password');
    if(password.value !== document.getElementById('repeat').value){
      this.addAlert("Invalide Password", "Passwords do no match");
    } else if(!password.value || password.value === '') {
      this.addAlert("Invalide Password", "You didn't write any new password");
    } else if(password.value.length < minPassLength) {
      this.addAlert("Invalide Password", "Password too short");
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
      <form >
        <ToastContainer toastMessageFactory={ToastMessageFactory} ref="container" className="toast-bottom-left" />

        <div className="row">
          {NodeInput}
        </div>

        <div className="row">
          <div className="form-group">

            <div className="col-sm-6">
              <label htmlFor="user-password">User Password </label>
              <input 
                type="password" className="form-control"
                id="password"
                placeholder="new password"
              />
            </div>

            <div className="col-sm-6">
              <label htmlFor="user-password">User Password </label>
              <input 
                type="password" className="form-control"
                id="repeat"
                placeholder="repeat password"
              />
            </div>

          </div>
        </div>

        <hr />

        <div className="clearfix">
          <div className="pull-right">

            <button className="btn btn-primary"
              onClick={this.updatePassword} >
              Update Password
            </button>

            <button className="btn btn-primary"
             onClick={this.handleSubmit} >
              <i className="fa fa-save"></i> Save
            </button>
          </div>
        </div>

      </form>
    );
  }
});

module.exports = UserForm;
