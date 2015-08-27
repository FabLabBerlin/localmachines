var {endDate, formatDate} = require('./helpers');
var moment = require('moment');
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
        var startDate = moment(membership.StartDate);
        var totalPrice = membership.Price;
        var priceExclVat = (Math.round(totalPrice / 1.19 * 100) / 100).toFixed(2);
        var vat = (totalPrice - priceExclVat).toFixed(2);
        return (
          <tr key={membership.Id} >
            <td>{membership.Title}</td>
            <td>{formatDate(startDate)}</td>
            <td>{formatDate(endDate(startDate, membership.Duration))}</td>
            <td>{priceExclVat} €</td>
            <td>{vat} €</td>
            <td>{totalPrice} €</td>
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
            <th>Start Date</th>
            <th>End Date</th>
            <th>Price/month excl. VAT</th>
            <th>VAT</th>
            <th>Total</th>
          </tr>
        </thead>
        <tbody>
          {MembershipNode}
        </tbody>
      </table>
    );
  }
});

export default Membership;
