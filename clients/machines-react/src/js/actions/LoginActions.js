import MachineStore from '../stores/MachineStore';

var LoginActions = {
    submitLoginForm: function(content) {
        MachineStore.submitLoginFormToServer(content);
    },

    logout() {
      MachineStore.logout();
    }
};

module.exports = LoginActions;
