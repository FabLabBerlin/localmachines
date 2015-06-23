var UserStore = {
    _state: {
        userID: 0,
        failToLogin: false,
        isLogged: false,
        // everything @ get /users/{uid}
        rawInfoUser: {
            Id: -1,
            FirstName: "fakeUser",
            LastName: "fakeUser",
            Username: "fakeUser",
            Email: "fakeUser@example.com",
            InvoiceAddr: 1,
            ShipAddr: 2,
            ClientId: 3,
            B2b: false,
            Company: "makea",
            VatUserId: "",
            VatRate: -1,
            UserRole: "",
            Created: "0001-01-01T00:00:20Z",
            Comments: ""    
        },
        // everything @ get /users/{uid}/machinepermissions
        rawInfoMachine: [
            {
                Id: '1',
                Name: 'ouioui',
                Shortname: 'oui',
                Description: 'du oui à foison'
            }
        ]
    },

   /*
    * To make request from the server
    * _url: the url for the API call
    * nameState: the name of the state you'll modify
    * _data: the data you're sending

    * _type: methods you use ('GET', 'POST' etc ...)
    *
    */
    getDataFromServer(_url, stateName, _dataToSend, _type, functionIfFail) {
        $.ajax({
            url: _url,
            dataType: 'json',
            type: _type,
            data: _dataToSend,
            success: function(data) {
                this._state[stateName] = data;
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(_url, status, err.toString());
                functionIfFail();
            }.bind(this),
        });
    },

    // It seems to work
    submitLoginFormToServer(loginInfo) {
        console.log(loginInfo);
        $.ajax({
            url: '/api/users/login',
            dataType: 'json',
            type: 'POST',
            data: loginInfo,
            success: function(data) {
                this._state.isLogged = true;
            }.bind(this),
            error: function(xhr, status, err) {
                console.error('/users/login', status, err.toString());
                this._state.failToLogin = true;
                //this.onChange();
            }.bind(this),
        });
    },

    // Logout from server
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
        /*
        var _url = '/api/users/' + this._state.userID,

            dataToSend = this._state.infoUser,
            _type = 'PUT';
        getDataFromServer(_url, dataToSend, _type, function() {});
        */
    },

    // To get the User information for the server after login in
    getUserStateFromServer(){
        //getting the initial set of data from the server with the UID
        //CALL API
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
        console.log(lightState);
        return lightState;
    },

    cleanState() {
        this._state.isLogged = false;
        this._state.failToLogin = false;
        this._state.userID = 0;
        this._state.rawInfoUser = {};
        this._state.rawInfoMachine = [];
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

    onChange() {}

};

module.exports = UserStore;
