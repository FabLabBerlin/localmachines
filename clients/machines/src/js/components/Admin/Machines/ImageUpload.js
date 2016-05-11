var _ = require('lodash');
var $ = require('jquery');
var Machines = require('../../../modules/Machines');
var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('../../../toastr');


var ImageUpload = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      newMachineImages: Machines.getters.getNewMachineImages
    };
  },

  fileChange: function(e) {
    const files = e.target.files;
    const mid = this.props.machine.get('Id');

    if (files) {
      var f = files[0];
      var reader = new window.FileReader();
      reader.onloadend = function() {
        console.log('calling upload action...');
        MachineActions.uploadMachineImage(mid, {
          dataUri: reader.result,
          fileName: f.name,
          fileSize: f.size
        });
      };
      reader.readAsDataURL(f);
    }
  },

  render() {
    const machine = this.props.machine;

    var machineImageFile;
    var machineImageNewFile;
    var machineImageNewFileName;
    var machineImageNewFileSize;

    if (machine.get('Image')) {
      machineImageFile = '/files/' + machine.get('Image');
    }

    if (this.state.newMachineImages.get(machine.get('Id'))) {
      machineImageNewFile = this.state.newMachineImages.get(machine.get('Id')).dataUri;
    }

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
          <div className="row"> 
            <div className="col-sm-6">
              <div className="form-group">
                <img id="machine-image"
                  src={machineImageNewFile || machineImageFile || 'assets/img/img-machine-placeholder.svg'}
                  alt="Machine image"/>
              </div>
            </div>
            <div className="col-sm-6">
              <div className="form-group">
                <input type="file"
                       onChange={this.fileChange}/>
                <button className="btn btn-primary btn-block"
                        disabled={!machineImageNewFile}
                        onClick={this.replaceImage}>
                  <i className="fa fa-file-image-o"/>&nbsp;Replace
                </button>
              </div>
              <p>Supported types: svg</p>
              <p>Name: {machineImageNewFileName || machine.get('ImageName')}<br/>Size: {machineImageNewFileSize || machine.get('ImageSize')}</p>
            </div>
          </div>

        </div>

      </div>
    );
  },

  replaceImage() {
    const mid = this.props.machine.get('Id');
    var data;

    if (this.state.newMachineImages.get(mid)) {
      data = this.state.newMachineImages.get(mid);
    }

    toastr.info('Uploading machine image...');

    $.ajax({
      method: 'POST',
      url: '/api/machines/' + mid + '/image',
      data: {
        Filename: data.fileName,
        Image: data.dataUri
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(){
      toastr.success('Machine image successfully uploaded');
    })
    .error(function(){
      toastr.error('Uploading machine image failed');
    });
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    console.log('updating', name, 'with', e.target.value);
    MachineActions.updateMachineField(id, name, e.target.value);
  }
});

export default ImageUpload;
