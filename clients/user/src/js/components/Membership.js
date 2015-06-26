import React from 'react';

var Membership = React.createClass({
  render() {
    var MembershipNode = ( 
      <tr>
        <td colSpan="2">You do not have any memberships</td>
      </tr> 
    );
    if(this.props.info.length != 0) {
      var membership = this.props.info[0];
      MembershipNode = (
        <tr>
          <td> Membership id: {membership.MembershipId}</td>
          <td>Start date: {membership.StartDate}</td>
        </tr>
      );
    }
    return (
      <table className="table table-striped table-hover">
        <thead>
          <tr>
            <th>Membership Id</th>
            <th>Start date</th>
          </tr>
        </thead>
        <tbody>
          {MembershipNode}
        </tbody>
      </table>
    );
  }
});

module.exports = Membership;
