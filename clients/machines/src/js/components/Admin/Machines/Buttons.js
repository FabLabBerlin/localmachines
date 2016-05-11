var $ = require('jquery');
var getters = require('../../../getters');
var LocationGetters = require('../../../modules/Location/getters');
var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('../../../toastr');


var Buttons = React.createClass({
  render() {
    const machine = this.props.machine;

    return (
      <div className="pull-right">

        {machine.get('Archived') ? (
          <button className="btn btn-danger"
                  onClick={this.toggleArchived}>
            <i className="fa fa-archive"></i>&nbsp;Unarchive
          </button>
        ) : (
          <button className="btn btn-danger"
                  onClick={this.toggleArchived}>
            <i className="fa fa-archive"></i>&nbsp;Archive
          </button>
        )}

        <button className="btn btn-primary"
                onClick={this.save}>
          <i className="fa fa-save"></i>&nbsp;Save
        </button>

      </div>
    );
  },

  save() {
    var machine = this.props.machine.toJS();

    console.log('machine:', machine);

    $.ajax({
      method: 'PUT',
      url: '/api/machines/' + machine.Id,
      contentType: 'application/json',
      data: JSON.stringify(machine),
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      toastr.success('Update successful');
    })
    .error(function(message, statusCode) {
      if (statusCode === 400 && message.indexOf('Dimensions') >= 0) {
        toastr.error(message);
      } else if (statusCode === 400 && message.indexOf('Found machine with same netswitch host') >= 0) {
        toastr.error(message);
      } else {
        toastr.error('Failed to update');
      }
    });
  },

  toggleArchived() {
    const machine = this.props.machine;
    const action = machine.get('Archived') ? 'unarchived' : 'archived';

    $.ajax({
      method: 'POST',
      url: '/api/machines/' + machine.get('Id') + '/set_archived?archived=' + !machine.get('Archived')
    })
    .success(function(data) {
      toastr.info('Successfully ' + action + ' machine');
      const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
      const uid = reactor.evaluateToJS(getters.getUid);
      MachineActions.apiGetUserMachines(locationId, uid);
    })
    .error(function() {
      toastr.error('Failed to ' + action + ' machine');
    });
  }
});

export default Buttons;
