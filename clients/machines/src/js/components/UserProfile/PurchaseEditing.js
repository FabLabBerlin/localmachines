var _ = require('lodash');
var {formatDuration} = require('./helpers');
var Invoices = require('../../modules/Invoices');
var React = require('react');


var Amount = React.createClass({
  render() {
    const p = this.props.purchase;

    if (p.Type === 'other') {
      return (
      <input type="number"
             onChange={this.update}
             value={this.props.purchase.Quantity}/>
      );
    } else {
      return <Duration invoice={this.props.invoice}
                       purchase={p}/>;
    }
  },

  update(e) {
    Invoices.actions.editPurchaseField({
      invoice: this.props.invoice,
      field: 'Quantity',
      value: e.target.value
    });
  }
});


var Category = React.createClass({
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
    Invoices.actions.editPurchaseField({
      invoice: this.props.invoice,
      field: 'Type',
      value: e.target.value
    });
  }
});


var Duration = React.createClass({
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


var PricePerUnit = React.createClass({
  render() {
    const p = this.props.purchase;

    if (p.Type !== 'other') {
      return <div>{p.PricePerUnit}</div>;
    }

    return (
      <input type="text"
             autoFocus="on"
             onChange={this.update}
             value={p.PricePerUnit}/>
    );
  },

  update(e) {
    Invoices.actions.editPurchaseField({
      invoice: this.props.invoice,
      field: 'PricePerUnit',
      value: e.target.value
    });
  }
});


var Unit = React.createClass({
  render() {
    const p = this.props.purchase;

    if (p.Type !== 'other') {
      return <div>{p.PriceUnit}</div>;
    }

    return (
      <select onChange={this.update}
              value={p.PriceUnit}>
        <option value="second">Second</option>
        <option value="minute">Minute</option>
        <option value="30 minutes">30 Minutes</option>
        <option value="hour">Hour</option>
        <option value="day">Day</option>
        <option value="gram">gram</option>
      </select>
    );
  },

  update(e) {
    Invoices.actions.editPurchaseField({
      invoice: this.props.invoice,
      field: 'PriceUnit',
      value: e.target.value
    });
  }
});

export default {
  Amount,
  Category,
  Duration,
  PricePerUnit,
  Unit
};
