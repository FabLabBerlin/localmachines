var UserStore = {
    _state: {
        userID: 0,
        failToLogin: false,
        isLogged: false,
        // everything @ get /users/{uid}
        rawInfoUser: {},
        // everything @ get /users/{uid}/machinepermissions
        rawInfoMachine: [],
        rawInfoMembership:{}
    },

    // Logout
    logoutFromServer() {
        $.ajax({
            url: '/api/users/logout',
            dataType: 'json',
            cache: false,
            success: function() {
                this.cleanState();
            }.bind(this),
            error: function(xhr, status, err) {
                console.error('/users/logout', status, err.toString())
            }.bind(this)
        });
    },

    // Put the information from UserForm in the state 
    // then send it to the server
    submitStateToServer(userState) {
        this.formatUserStateToSendToServer(userState);
    },

    // It seems to work
    submitLoginFormToServer(loginInfo) {
        $.ajax({
            url: '/api/users/login',
            dataType: 'json',
            type: 'POST',
            data: loginInfo,
            success: function(data) {
                this._state.UserID = data["UserId"];
                this.getUserStateFromServer(data["UserId"]);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error('/users/login', status, err.toString());
                this._state.failToLogin = true;
                //this.onChange();
            }.bind(this),
        });
    },

    // To get the User information from the server after login
    getUserStateFromServer(uid){
        $.ajax({
            url: '/api/users/' + uid,
            dataType: 'json',
            type: 'GET',
            success: function(data) {
                this._state.rawInfoUser = data;
                this.getMachineFromServer(uid);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error('/users/{uid}', status, err.toString());
            }.bind(this),
        });
    },

    // To get the Machine information from the server after login
    getMachineFromServer(uid){
        $.ajax({
            url: '/api/users/' + uid + '/machinepermissions',
            dataType: 'json',
            type: 'GET',
            success: function(data) {
                this._state.rawInfoMachine = data;
                this.getMembershipFromServer(uid);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error('/users/{uid}/machinepermissions', status, err.toString());
            }.bind(this),
        });
    },

    getMembershipFromServer(uid) {
        $.ajax({
            url: '/api/users/' + uid + '/memberships',
            dataType: 'json',
            type: 'GET',
            success: function(data) {
                this._state.rawInfoMembership = data;
                this._state.isLogged = true;
                this.onChange();
            }.bind(this),
            error: function(xhr, status, err) {
                console.error('/users/{uid}/memberships', status, err.toString());
            }.bind(this),
        });
    },

    // Change the input to match the rawInfoUser format
    formatUserStateToSendToServer(userState) {
        for(var data in userState) {
            this._state.rawInfoUser[data] = userState[data];
        }
    },

    // return a light state with only the useful information
    formatUserStateToSendToUserPage() {
        var infoWhichMatter = [
            'FirstName',
            'LastName', 
            'Username',
            'Email',
            'InvoiceAddr',
            'ShipAddr'
        ];
        var lightState = {};
        for(var index in infoWhichMatter) {
            var name = infoWhichMatter[index];
            lightState[name] = this._state.rawInfoUser[name];
        }
        return lightState;
    },

    // To call before logout
    cleanState() {
        this._state.isLogged = false;
        this._state.failToLogin = false;
        this._state.userID = 0;
        this._state.rawInfoUser = {};
        this._state.rawInfoMachine = [];
        this._state.rawInfoMachine = {};
    },

    // Getter to the state
    getIsLogged: function() {
        return this._state.isLogged;
    },
    // Use by UserPage to get its state
    getInfoUser: function() {
        var lightState = this.formatUserStateToSendToUserPage();
        return lightState;
    },
    // Use by UserPage to get its state
    getInfoMachine: function() {
        return this._state.rawInfoMachine;
    },

    getMembership() {
        return this._state.rawInfoMembership;
    },

    onChange() {}

};

module.exports = UserStore;
