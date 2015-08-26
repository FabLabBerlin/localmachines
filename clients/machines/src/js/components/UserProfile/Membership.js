import React from 'react';

/*
 * Membership component:
 * Display the membership the user is subscribing
 */
var Membership = React.createClass({

  /*
   * Create the table for each data
   * Display it
   */
  render() {
    var MembershipNode;
    if (this.props.info && this.props.info.length !== 0) {
      MembershipNode = this.props.info.map(function(membership) {
        return (
        <tr key={membership.Id} >
            <td>{membership.Title}</td>
            <td>{membership.Duration + ' ' + membership.Unit}</td>
            <td>{membership.Price}</td>
            <td>{membership.StartDate}</td>
          </tr>
        );
      });
    } else {
      return <p>You do not have any memberships</p>;
    }
    return (
      <table className="table table-striped table-hover">
        <thead>
          <tr>
            <th>Name</th>
            <th>Duration</th>
            <th>Price</th>
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
