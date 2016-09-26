var ImageUploader = require('../ImageUploader');
var Machines = require('../../../modules/Machines');
var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../reactor');


var ImageUpload = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      newMachineImages: Machines.getters.getNewMachineImages
    };
  },

  render() {
    const machine = this.props.machine;
    const machineImageFile = machine.get('Image') ? ('/files/' + machine.get('Image')) : null;

    return (
      <div className="row">

        <div className="col-sm-3">
          <div className="form-group">
            <label>Machine Description</label>
            <textarea className="form-control"
                      onChange={this.update.bind(this, 'Description')}
                      placeholder="Enter machine description"
                      value={machine.get('Description')}
                      rows="5" />
          </div>
        </div>

        <div className="col-sm-6">
          <label>Image</label>
          <ImageUploader existingImage={machineImageFile}
                         uploadUrl={'/api/machines/' + machine.get('Id') + '/image'}/>
        </div>

      </div>
    );
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    MachineActions.updateMachineField(id, name, e.target.value);
  }
});

export default ImageUpload;
