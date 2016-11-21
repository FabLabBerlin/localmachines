var _ = require('lodash');
var {formatDuration} = require('./helpers');
var getters = require('../../getters');
var Invoices = require('../../modules/Invoices');
var LoaderLocal = require('../LoaderLocal');
var LocationGetters = require('../../modules/Location/getters');
var MachineActions = require('../../actions/MachineActions');
var Machines = require('../../modules/Machines');
var moment = require('moment');
var React = require('react');
var Datetime = require('react-datetime');
var reactor = require('../../reactor');
var UserActions = require('../../actions/UserActions');


var Amount = React.createClass({
  render() {
    const p = this.props.purchase;

    if (p.PriceUnit === 'gram') {
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
    if (!this.props.purchase.PriceUnit) {
      return <div/>;
    }

    return (
      <input type="text"
             autoFocus="on"
             ref="duration"
             onChange={this.update}
             placeholder="0:00"
             value={this.props.purchase.editedDuration ||
                    formatDuration(this.props.purchase)}/>
    );
  },

  update(e) {
    Invoices.actions.editPurchaseDuration(this.props.invoice, e.target.value);
  }
});


var Name = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocation).Id;
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    MachineActions.apiGetUserMachines(locationId, uid);
  },

  render() {
    const p = this.props.purchase;

    if (p.Type === 'other') {
      return (
        <input type="text"
               onChange={this.update}
               value={p.CustomName}/>
      );
    } else {
      if (!this.state.machines) {
        return <LoaderLocal/>;
      }

      return (
        <select onChange={this.update}
                value={p.MachineId}>
          <option value="0">Please select</option>
          {this.state.machines.toList()
                              .sortBy(m => m.get('Name'))
                              .map(m => {
            return (
              <option key={m.get('Id')}
                      value={m.get('Id')}>
                {m.get('Name')}
              </option>
            );
          })}
        </select>
      );
    }
  },

  update(e) {
    const p = this.props.purchase;

    if (p.Type === 'other') {
      Invoices.actions.editPurchaseField({
        invoice: this.props.invoice,
        field: 'CustomName',
        value: e.target.value
      });
    } else {
      Invoices.actions.editPurchaseField({
        invoice: this.props.invoice,
        field: 'MachineId',
        value: e.target.value
      });
    }
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


var StartTime = React.createClass({
  render() {
    const p = this.props.purchase;

    return <Datetime value={moment(p.TimeStart)}
                     onChange={this.update}
                     dateFormat="DD. MMM YYYY"
                     timeFormat="HH:mm"/>;
  },

  update(t) {
    Invoices.actions.editPurchaseField({
      invoice: this.props.invoice,
      field: 'TimeStart',
      value: t.toJSON()
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
  Name,
  PricePerUnit,
  StartTime,
  Unit
};
