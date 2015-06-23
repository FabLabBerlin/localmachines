import UserStore from '../stores/UserStore'

var UserActions = {

    submitState(userState){
        UserStore.submitStateToServer(userState);
    },

    logout() {
        UserStore.logoutFromServer();
    }

};

module.exports = UserActions;
