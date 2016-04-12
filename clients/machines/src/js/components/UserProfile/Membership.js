var {formatDate, subtractVAT, toCents, toEuro} = require('./helpers');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var SettingsGetters = require('../../modules/Settings/getters');

/*
 * Membership component:
 * Display the membership the user is subscribing
 */
var Membership = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      vatPercent: SettingsGetters.getVatPercent
    };
  },

  /*
   * Create the table for each data
   * Display it
   */
  render() {
    var MembershipNode;
    if (this.props.memberships && this.props.memberships.length !== 0) {
      MembershipNode = this.props.memberships.map(function(membership) {
        var startDate = moment(membership.StartDate);
        var endDate = moment(membership.EndDate);
        var totalPrice = toCents(membership.MonthlyPrice);
        var priceExclVat = toCents(subtractVAT(membership.MonthlyPrice));
        var vat = totalPrice - priceExclVat;
        return (
          <tr key={membership.Id} >
            <td>{membership.Title}</td>
            <td>{formatDate(startDate)}</td>
            <td>{formatDate(endDate)}</td>
            <td>{toEuro(priceExclVat)}€</td>
            <td>{toEuro(vat)}€</td>
            <td>{toEuro(totalPrice)}€</td>
          </tr>
        );
      });
    } else {
      return <p>You do not have any memberships</p>;
    }
    return (
      <div className="table-responsive">
        <table className="table table-stripped table-hover">
          <thead>
            <tr>
              <th>Name</th>
              <th>Start Date</th>
              <th>End Date</th>
              <th>Price/month excl. VAT</th>
              <th>VAT/month ({this.state.vatPercent})</th>
              <th>Total/month</th>
            </tr>
          </thead>
          <tbody>
            {MembershipNode}
          </tbody>
        </table>
      </div>
    );
  }
});

export default Membership;
