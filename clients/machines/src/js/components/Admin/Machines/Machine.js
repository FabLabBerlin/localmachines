import _ from 'lodash';
import BasicData from './BasicData';
import Buttons from './Buttons';
import getters from '../../../getters';
import ImageUpload from './ImageUpload';
import LoaderLocal from '../../LoaderLocal';
import LocationGetters from '../../../modules/Location/getters';
import Machines from '../../../modules/Machines';
import MachineActions from '../../../actions/MachineActions';
import MachineProperties from './MachineProperties';
import NetswitchConfig from './NetswitchConfig';
import React from 'react';
import reactor from '../../../reactor';
import toastr from '../../../toastr';


var Machine = React.createClass({

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
    const machineId = parseInt(this.props.params.machineId);
    var machine;
    if (this.state.machines) {
      machine = this.state.machines.find((m) => {
        return m.get('Id') === machineId;
      });
    }
    if (machine) {
      return (
        <div className="container-fluid">
          <h1>Edit Machine</h1>

          <hr />

          <BasicData machine={machine} />
          <ImageUpload machine={machine} />
          <MachineProperties machine={machine} />
          <NetswitchConfig machine={machine} />

          <Buttons machine={machine} />
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }

});

export default Machine;
