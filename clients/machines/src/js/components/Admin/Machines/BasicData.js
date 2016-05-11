var MachineActions = require('../../../actions/MachineActions');
var React = require('react');


var BasicData = React.createClass({
  render() {
    const machine = this.props.machine;

    return (
      <div className="row">

        <div className="col-sm-3">
          <div className="form-group">
            <label>Machine Name</label>
            <input id="machine-name"
                   type="text"
                   className="form-control"
                   onChange={this.update.bind(this, 'Name')}
                   placeholder="Enter machine name"
                   value={machine.get('Name')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Short Name</label>
            <input type="text"
                   onChange={this.update.bind(this, 'Shortname')}
                   className="form-control"
                   placeholder="Enter short name"
                   value={machine.get('Shortname')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Price</label>
            <input type="text"
                   onChange={this.update.bind(this, 'Price')}
                   className="form-control"
                   placeholder="Enter price"
                   value={machine.get('Price')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Price Unit</label>
            <select className="form-control"
                    onChange={this.update.bind(this, 'PriceUnit')}
                    id="machine-price-unit"
                    value={machine.get('PriceUnit')}>
              <option value="" disabled>Select unit</option>
              <option value="minute">minute</option>
              <option value="hour">hour</option>
              <option value="day">day</option>
            </select>
          </div>
        </div>

      </div>
    );
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    var value = e.target.value;
    switch (name) {
    case 'Price':
      value = parseFloat(value);
      break;
    }
    MachineActions.updateMachineField(id, name, value);
  }
});

export default BasicData;
