var LoaderLocal = require('../LoaderLocal');
var LocationGetters = require('../../modules/Location/getters');
var MachineActions = require('../../actions/MachineActions');
var Machines = require('../../modules/Machines');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


var Machine = React.createClass({
  imgUrl() {
    if (this.props.machine.get('Image')) {
      return '/files/' + this.props.machine.get('Image');
    } else {
      return '/machines/img/img-machine-placeholder.svg';
    }
  },

  render() {
    const m = this.props.machine;
    const style = {
      backgroundImage: 'url(' + this.imgUrl() + ')'
    };

    return (
      <div className="ms-machine">
        <div className="ms-machine-label">
          {m.get('Name')}
        </div>
        <div className="ms-machine-icon" style={style}>
        </div>
      </div>
    );
  }
});


var MachinesPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      locationId: LocationGetters.getLocationId,
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);

    if (window.WebSocket) {
      MachineActions.wsDashboard(this.context.router, locationId);
    } else {
      MachineActions.lpDashboard(this.context.router, locationId);
    }
  },

  render() {
    if (!this.state.locationId || !this.state.machines) {
      return <LoaderLocal/>;
    }

    var machines = this.state.machines
    .toList()
    .filter(m => {
      return m.get('LocationId') === this.state.locationId &&
        !m.get('Archived');
    })
    .sortBy(m => m.get('Name'));

    return (
      <div id="ms" className="container-fluid">
        {machines.map((m, i) => <Machine key={i} machine={m}/>)}
      </div>
    );
  }

});

export default MachinesPage;
