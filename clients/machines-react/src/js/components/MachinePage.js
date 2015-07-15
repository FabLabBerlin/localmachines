import React from 'react'
import MachineList from './MachineList';
import MachineStore from '../stores/MachineStore';

var MachinePage = React.createClass({

  getInitialState() {
    return {
      userInfo: MachineStore.getUserInfo(),
      machineInfo: MachineStore.getMachineInfo(),
      activationInfo: MachineStore.getActivationInfo()
    };
  },

  getUserId() {
    return this.state.userInfo.uid;
  },

  onChangeActivation() {
    this.setState({
      activationInfo: MachineStore.getActivationInfo()
    });
  },

  componentDidMount() {
    MachineStore.onChangeActivation = this.onChangeActivation;
  },

  render() {
    return (
      <div className="container-fluid" >
        <div>
          coucou {this.state.userInfo.FirstName} {this.state.userInfo.LastName}
        </div>
        <div>
          <MachineList 
            uid={this.getUserId()}
            info={this.state.machineInfo} 
            activation={this.state.activationInfo}
          />
        </div>
      </div>
    );
  }
});

module.exports = MachinePage;
