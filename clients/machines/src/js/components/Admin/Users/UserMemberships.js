var DatePicker = require('./DatePicker');
var LocationGetters = require('../../../modules/Location/getters');
var React = require('react');
var reactor = require('../../../reactor');


var UserMemberships = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    };
  },

  render() {
    var membership = {};
    var userMembership = {};

    return (
      <div>
        <h2>User Memberships</h2>
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Membership Name</th>
              <th>Start Date</th>
              <th>End Date</th>
              <th>Extends Automatically</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr ng-hide="userMemberships.length">
              <td colspan="6">User has no memberships</td>
            </tr>
            <tr ng-repeat="userMembership in userMemberships">
              <td>{userMembership.Title}</td>
              <td style={{minWidth: '90px'}}>{userMembership.StartDate}</td>
              <td style={{minWidth: '160px'}}>
                <div className="form-inline">
                  <DatePicker placeholder="End Date"
                              value={userMembership.EndDate}/>
                </div>
              </td>
              <td>
                <input 
                  type="checkbox" 
                  ng-model="userMembership.AutoExtend"
                  ng-disabled="userMemberhip.Inactive"
                  ng-change="updateUserMembership(userMembership.Id)"/>
              </td>
              <td>
                <div ng-show="userMembership.Active">Active</div>
                <div ng-show="userMembership.Cancelled">Cancelled</div>
                <div ng-show="userMembership.Inactive">Inactive</div>
              </td>
            </tr>
          </tbody>
        </table>

        <div className="row">
          <div className="col-sm-3">
            <div className="form-group">
              <select 
                className="form-control" 
                id="user-select-membership" 
                placeholder="Membership">
                <option value="" selected disabled>Select Membership</option>
                <option ng-repeat="membership in memberships | membershipsFilter:this" 
                        value={membership.Id}>
                  {membership.Title}
                </option>
              </select>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <DatePicker placeholder="Start Date"/>
            </div>
          </div>

          <div className="col-sm-3">
            <button 
              className="btn btn-primary btn-block" 
              id="user-add-membership-btn" 
              ng-click="addUserMembership()">
              <i className="fa fa-plus"></i>&nbsp;Add Membership
            </button>
          </div>
        </div>
      </div>
    );
  }

});

export default UserMemberships;
