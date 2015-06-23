import UserStore from '../stores/UserStore';

var LoginActions = {
    submitLoginForm: function(content) {
        UserStore.submitLoginFormToServer(content);
    }
};

module.exports = LoginActions;
