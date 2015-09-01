var {endDate, formatDate, subtractVAT, toCents, toEuro} = require('./helpers');
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
        var totalPrice = toCents(membership.MonthlyPrice);
        var priceExclVat = toCents(subtractVAT(membership.MonthlyPrice));
        var vat = totalPrice - priceExclVat;
        return (
          <tr key={membership.Id} >
            <td>{membership.Title}</td>
            <td>{formatDate(startDate)}</td>
            <td>{formatDate(endDate(startDate, membership.Duration))}</td>
            <td>{toEuro(priceExclVat)} €</td>
            <td>{toEuro(vat)} €</td>
            <td>{toEuro(totalPrice)} €</td>
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
            <th>VAT (19%)</th>
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
