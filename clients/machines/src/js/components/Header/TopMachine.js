var Machines = require('../../modules/Machines');
var React = require('react');
var reactor = require('../../reactor');


var TopMachine = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machines: Machines.getters.getMachines
    };
  },

  hide() {
    window.history.back();
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
    console.log('machineId=', this.props.machineId);
    const m = this.machine();

    return (
      <div className="nav-top row">
        {m ? (
          <div className="nav-machine-top-panel">
            <div className="nav-machine-name">{m.get('Name')}</div>
            <div></div>
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
