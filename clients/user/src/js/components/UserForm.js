import _ from 'lodash';
import React from 'react';

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
      <form onSubmit={this.handleSubmit}>

        <div className="row">
          {NodeInput}
        </div>

        <div className="row">
          <div className="form-group">

            <div className="col-sm-6">
              <label htmlFor="user-password">User Password </label>
              <input 
                type="password" className="form-control"
                placeholder="new password"
              />
            </div>

            <div className="col-sm-6">
              <label htmlFor="user-password">User Password </label>
              <input 
                type="password" className="form-control"
                placeholder="repeat password"
              />
            </div>

          </div>
        </div>

        <hr />

        <div className="clearfix">
          <div className="pull-right">
            <button className="btn btn-primary">
              <i className="fa fa-save"></i> Save
            </button>
          </div>
        </div>

      </form>
    );
  }
});

module.exports = UserForm;
