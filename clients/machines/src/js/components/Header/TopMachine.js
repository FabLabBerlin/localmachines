import {hashHistory} from 'react-router';
import Machines from '../../modules/Machines';
import React from 'react';
import reactor from '../../reactor';
import Settings from '../../modules/Settings';


var MachineType = React.createClass({
  render() {
    const tId = this.props.machine.get('TypeId');
    const t = {
      1: '3D Printer',
      2: 'CNC Mill',
      3: 'Heatpress',
      4: 'Knitting Machine',
      5: 'Laser Cutter',
      6: 'Vinylcutter'
    }[tId];

    return (
      <div className="nav-machine-type">
        {t}
      </div>
    );
  }
});


var Price = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      currency: Settings.getters.getCurrency
    };
  },

  render() {
    var price;
    price = this.props.machine.get('Price').toFixed(2);
    price += ' ' + (this.state.currency || '€') + ' / ';
    switch (this.props.machine.get('PriceUnit')) {
      case 'hour':
        price += 'h';
        break;
      case 'minute':
        price += 'min';
        break;
      default:
        price += this.props.machine.get('PriceUnit');
    }

    return (
      <div>
        {price} (incl. VAT)
      </div>
    );
  }
});


var TopMachine = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machines: Machines.getters.getMachines
    };
  },

  hide() {
    hashHistory.push('/machines');
  },

  machine() {
    var m;

    if (this.state.machines) {
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === this.props.machineId;
      });
    }

    return m;
  },

  render() {
    const m = this.machine();

    return (
      <div className="nav-top row">
        {m ? (
          <div className="nav-machine-top-panel">
            <div className="nav-machine-info">
              <div className="nav-machine-name">{m.get('Name')}</div>
              <MachineType machine={m}/>
              <Price machine={m}/>
            </div>
            <button type="button"
                    title="Close"
                    onClick={this.hide}>
              <img src="/machines/assets/img/machines/machine/m_cancel.svg"/>
            </button>
          </div>
        ) : null}
      </div>
    );
  }
});

export default TopMachine;
