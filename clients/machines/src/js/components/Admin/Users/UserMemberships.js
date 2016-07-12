var _ = require('lodash');
var $ = require('jquery');
var DatePicker = require('react-datepicker');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');
var Users = require('../../../modules/Users');

require('react-datepicker/dist/react-datepicker.css');


var UserMembership = React.createClass({
  render() {
    const userMembership = this.props.userMembership;
    const endDate = moment(userMembership.EndDate);

    return (
      <tr>
        <td>{userMembership.Title}</td>
        <td style={{minWidth: 90}}>{userMembership.StartDate}</td>
        <td style={{minWidth: 160}}>
          <div className="form-inline">
            <DatePicker dateFormat="YYYY-MM-DD"
                        placeholder="End Date"
                        selected={endDate}/>
          </div>
        </td>
        <td>
          <input 
            type="checkbox" 
            value={userMembership.AutoExtend}
            disabled={userMembership.Inactive}
            ng-change="updateUserMembership(userMembership.Id)"/>
        </td>
        <td>
          <div ng-show="userMembership.Active">Active</div>
          <div ng-show="userMembership.Cancelled">Cancelled</div>
          <div ng-show="userMembership.Inactive">Inactive</div>
        </td>
      </tr>
    );
  }
});


var UserMemberships = React.createClass({

  mixins: [ reactor.ReactMixin ],

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const userId = reactor.evaluateToJS(getters.getUid);

    Users.actions.fetchMemberships({locationId});
    Users.actions.fetchUserMemberships({locationId, userId});
  },

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      memberships: Users.getters.getMemberships,
      userMemberships: Users.getters.getUserMemberships
    };
  },

  add() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const userId = reactor.evaluateToJS(getters.getUid);
    const membershipId = parseInt(
      $(this.refs.newMembershipId.getDOMNode()).val()
    );

    Users.actions.addUserMembership({locationId, userId, membershipId});
  },

  render() {
    const memberships = this.state.memberships;
    const user = this.props.user;
    console.log('user=', user);
    const userMemberships = this.state.userMemberships.get(user.Id);

    if (!userMemberships) {
      return <LoaderLocal/>;
    }

    return (
      <div>
        <h2>User Memberships</h2>
        {userMemberships.length ? (
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
              {_.map(userMemberships, (userMembership, i) => {
                return <UserMembership key={i}
                                       userMembership={userMembership}/>;
              })}
            </tbody>
          </table>
        ) : (
          <div>User has no memberships</div>
        )}

        {memberships ? (
          <div className="row">
            <div className="col-sm-3">
              <div className="form-group">
                <select className="form-control" 
                        id="user-select-membership" 
                        placeholder="Membership"
                        ref="newMembershipId">
                  <option value="" selected disabled>Select Membership</option>
                  {_.map(memberships, (membership) => {
                    return (
                      <option value={membership.Id}>
                        {membership.Title}
                      </option>
                    );
                  })}
                </select>
              </div>
            </div>

            <div className="col-sm-3">
              <button 
                className="btn btn-primary btn-block" 
                id="user-add-membership-btn" 
                onClick={this.add}>
                <i className="fa fa-plus"></i>&nbsp;Add Membership
              </button>
            </div>
          </div>
        ) : <LoaderLocal/>}

      </div>
    );
  }

});

export default UserMemberships;
