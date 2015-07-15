import MachineStore from '../stores/MachineStore';

var LoginActions = {
    submitLoginForm: function(content) {
        MachineStore.submitLoginFormToServer(content);
    }
};

module.exports = LoginActions;
