import MachineStore from '../stores/MachineStore';

var LoginActions = {
    submitLoginForm: function(content) {
        MachineStore.apiPostLogin(content);
    },

    logout() {
      MachineStore.apiGetLogout();
    }
};

module.exports = LoginActions;
