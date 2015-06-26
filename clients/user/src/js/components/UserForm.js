import _ from 'lodash';
import React from 'react';

var UserForm = React.createClass({

  // Pass changing of the input responsability to the parent
  handleChangeForm(event) {
    this.props.func(event);
  },

  // Pass the submit responsibility to the parent
  handleSubmit() {
    // should send value to useractions
    this.props.submit();
  },

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
      //l61, put font-awesome to get the logo
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
