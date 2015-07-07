import UserStore from '../stores/UserStore'

var UserActions = {

    submitState(userState){
        UserStore.submitUpdatedStateToServer(userState);
    },

    logout() {
        UserStore.logoutFromServer();
    }

};

module.exports = UserActions;
