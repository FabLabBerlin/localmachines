import _ from 'lodash';
import getters from '../../../getters';
import LocationGetters from '../../../modules/Location/getters';
import MachineActions from '../../../actions/MachineActions';
import Machines from '../../../modules/Machines';
import React from 'react';
import reactor from '../../../reactor';


var MachinesView = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
  },

  render() {
    var machines = _.sortBy(this.state.machines.toJS(), m => m.Name);
    return (
      <div className="container-fluid">
        <div className="row">
          <div className="col-xs-1">
          </div>
          <div className="col-xs-11">
            {_.map(machines, (m) => {
              return (
                <div key={m.Id}>
                  <div className="col-xs-6">
                    <div className="row">
                      {m.Name}
                    </div>
                  </div>
                  <div className="col-xs-5 text-center">
                    <a type="button"
                       className="btn btn-primary btn-ico pull-right"
                       href={'/machines/#/admin/machines/' + m.Id}>
                      <i className="fa fa-edit"></i>
                    </a>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    );
  }

});

export default MachinesView;
