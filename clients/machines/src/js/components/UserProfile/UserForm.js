import _ from 'lodash';
import countryCodes from './CountryCodes';
import React from 'react';
import toastr from '../../toastr';


var LabelledText = React.createClass({
  render() {
    return (
      <div className="col-md-6">
        <div className="form-group">
          <label htmlFor={this.props.id}>{this.props.label}</label>
          <input
            type="text"
            value={this.props.value}
            id={this.props.id}
            className="form-control"
            onChange={this.props.onChange}
          />
        </div>
      </div>
    );    
  }
});


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

  render() {
    return (
      <form >

        <div className="row">
          {this.props.user.get('Company') ?
            <LabelledText id="Company" label="Company"
                          value={this.props.user.get('Company')}
                          onChange={this.handleChangeForm}/> : null}
          <LabelledText id="FirstName" label="First Name"
                        value={this.props.user.get('FirstName')}
                        onChange={this.handleChangeForm}/>
          <LabelledText id="LastName" label="Last Name"
                        value={this.props.user.get('LastName')}
                        onChange={this.handleChangeForm}/>
          <LabelledText id="Username" label="Username"
                        value={this.props.user.get('Username')}
                        onChange={this.handleChangeForm}/>
          <LabelledText id="Email" label="E-Mail"
                        value={this.props.user.get('Email')}
                        onChange={this.handleChangeForm}/>
          <LabelledText id="Phone" label="Phone"
                        value={this.props.user.get('Phone')}
                        onChange={this.handleChangeForm}/>
          <LabelledText id="InvoiceAddr" label="Invoice Address"
                        value={this.props.user.get('InvoiceAddr')}
                        onChange={this.handleChangeForm}/>
          <LabelledText id="ZipCode" label="Zip Code"
                        value={this.props.user.get('ZipCode')}
                        onChange={this.handleChangeForm}/>
          <LabelledText id="City" label="City"
                        value={this.props.user.get('City')}
                        onChange={this.handleChangeForm}/>
          <div className="col-md-6">
            <div className="form-group">
              <label htmlFor="CountryCode">Country</label>
              <select id="CountryCode"
                      className="form-control"
                      defaultValue={this.props.user.get('CountryCode')}
                      onChange={this.handleChangeForm}>
                <option value="" disabled>Select Country</option>
                {_.map(countryCodes, function(value, key) {
                  return (
                    <option value={value.Code}>{value.Name}</option>
                  );
                })}
              </select>
            </div>
          </div>

          <div className="col-md-12" >
            <label htmlFor="user-information" >Password</label>
          </div>

          <div className="form-group">
            <div className="col-sm-3">
              <input
                type="password" className="form-control"
                id="password"
                placeholder="new password"
              />
            </div>
            <div className="col-sm-3">
              <input
                type="password" className="form-control"
                id="repeat"
                placeholder="repeat password"
              />
            </div>
            <div className="col-sm-2" >
              <button className="btn btn-primary"
                onClick={this.updatePassword} >
                Update Password
              </button>
            </div>
          </div>
        </div>


        <hr />

        <div className="clearfix">
          <div className="pull-right">

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

export default UserForm;
