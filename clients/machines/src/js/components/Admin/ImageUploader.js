var $ = require('jquery');
import React from 'react';
import toastr from '../../toastr';


var ImageUploader = React.createClass({
  getInitialState() {
    return {
      newImage: undefined
    };
  },

  fileChange: function(e) {
    const files = e.target.files;

    if (files) {
      var f = files[0];
      var reader = new window.FileReader();
      reader.onloadend = () => {
        this.setState({
          newImage: {
            dataUri: reader.result,
            fileName: f.name,
            fileSize: f.size
          }
        });
      };
      reader.readAsDataURL(f);
    }
  },

  render() {
    var machineImageFile;
    var machineImageNewFile;
    var machineImageNewFileName;
    var machineImageNewFileSize;

    if (this.state.newImage) {
      machineImageNewFile = this.state.newImage.dataUri;
      machineImageNewFileName = this.state.newImage.fileName;
    }

    return (
      <div className="row"> 
        <div className="col-sm-6">
          <div className="form-group">
            <img id="machine-image"
                 height={this.props.height}
                 src={machineImageNewFile || this.props.existingImage || 'assets/img/img-machine-placeholder.svg'}
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
          <p>{machineImageNewFileSize ? ('Size: ' + machineImageNewFileSize) : ''}</p>
        </div>
      </div>
    );
  },

  replaceImage() {
    var data;

    if (this.state.newImage) {
      data = this.state.newImage;
    } else {
      return;
    }

    toastr.info('Uploading image...');

    $.ajax({
      method: 'POST',
      url: this.props.uploadUrl,
      data: {
        Filename: data.fileName,
        Image: data.dataUri
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(){
      toastr.success('Image successfully uploaded');
    })
    .error(function(){
      toastr.error('Uploading image failed');
    });
  }
});

export default ImageUploader;
