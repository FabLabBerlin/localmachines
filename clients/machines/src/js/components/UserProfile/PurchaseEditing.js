var {formatDuration} = require('./helpers');
var Invoices = require('../../modules/Invoices');
var React = require('react');


var CategoryEdit = React.createClass({
  render() {
    const p = this.props.purchase;

    if (p.Type === 'tutor') {
      return <div>Tutoring</div>;
    }

    return (
      <select onChange={this.update}
              value={p.Type}>
        <option value="activation">Activation</option>
        <option value="reservation">Reservation</option>
        <option value="other">Other</option>
      </select>
    );
  },

  update(e) {
    Invoices.actions.editPurchaseCategory(this.props.invoice, e.target.value);
  }
});


var DurationEdit = React.createClass({
  render() {
    return (
      <input type="text"
             autoFocus="on"
             ref="duration"
             onChange={this.update}
             value={this.props.purchase.editedDuration ||
                    formatDuration(this.props.purchase)}/>
    );
  },

  update(e) {
    Invoices.actions.editPurchaseDuration(this.props.invoice, e.target.value);
  }
});


var UnitEdit = React.createClass({
  render() {
    return (
      <input type="text"
             autoFocus="on"
             onChange={this.update}
             value={this.props.purchase.PriceUnit}/>
    );
  },

  update(e) {
    Invoices.actions.editPurchaseUnit(this.props.invoice, e.target.value);
  }
});

export default {
  CategoryEdit,
  DurationEdit,
  UnitEdit
};
