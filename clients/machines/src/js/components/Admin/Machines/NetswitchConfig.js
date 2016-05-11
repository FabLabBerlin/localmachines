var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var toastr = require('../../../toastr');


var NetswitchConfig = React.createClass({
  applyConfig() {
    toastr.error('Not implemented here yet');
  },

  render() {
    const machine = this.props.machine;

    var netswitchConfigStatus;

    return (
      <div className="row">

        <div className="col-sm-12">
          <label>NetSwitch Config</label>
          <div className="row">

            {machine.get('NetswitchType') ?
              (
                <div className="col-sm-3">
                  <div className="form-group">
                    <input type="text"
                           className="form-control"
                           placeholder="Host (mfi only)"
                           value={machine.get('NetswitchHost')}
                           onChange={this.update.bind(this, 'NetswitchHost')}/>
                  </div>
                </div>
              ) : (
                <div>
                  <div className="col-sm-3">
                    <div className="form-group">
                      <input type="text"
                             className="form-control"
                             placeholder="On URL"
                             value={machine.get('NetswitchUrlOn')}
                             onChange={this.update.bind(this, 'NetswitchUrlOn')}/>
                    </div>
                  </div>
                  <div className="col-sm-3">
                    <div className="form-group">
                      <input type="text"
                             className="form-control"
                             placeholder="Off URL"
                             value={machine.get('NetswitchUrlOff')}
                             onChange={this.update.bind(this, 'NetswitchUrlOff')}/>
                    </div>
                  </div>
                </div>
              )
            }

          </div>
          <div className="row">

            <div className="col-sm-3">
              <div className="form-group">
                <select className="form-control"
                        value={machine.get('NetswitchType')}
                        onChange={this.update.bind(this, 'NetswitchType')}>
                  <option value="">Custom Powerswitch</option>
                  <option value="mfi">Ubiquiti mFi Powerswitch</option>
                </select>
              </div>
            </div>

            <div className="col-sm-3">
              <div className="form-group">
                {netswitchConfigStatus ?
                  (
                    <button className="btn btn-danger btn-block"
                            disabled="true"
                            type="button">
                      <i className="fa fa-refresh fa-spin"></i>
                      {netswitchConfigStatus}
                    </button>
                  ) : null
                }
                {(!netswitchConfigStatus && machine.NetswitchType === 'mfi') ?
                  (
                    <button className="btn btn-danger btn-block"
                            onClick={this.applyConfig}
                            type="button">
                      Upgrade Powerswitch
                    </button>
                  ) : null
                }
              </div>
            </div>
            
          </div>
        </div>

      </div>
    );
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    console.log('updating', name, 'with', e.target.value);
    MachineActions.updateMachineField(id, name, e.target.value);
  }
});

export default NetswitchConfig;
