import MachineStore from '../stores/MachineStore';

var LoginActions = {
    submitLoginForm: function(content) {
        MachineStore.postLogin(content);
    },

    logout() {
      MachineStore.getLogout();
    }
};

module.exports = LoginActions;
