import React from 'react';

var UserForm = React.createClass({

    // Pass changing of the input responsability to the parent
    handleChangeForm(event) {
        this.props.func(event);
    },

    // Pass the submit responsability to the parent
    handleSubmit() {
        // should sent value to useractions
        this.props.submit();
    },

    render() {
        // Create a temporary Array to map it easily
        // Putting in the name(key) of the data(value) and his value
        var infoUserTmp = new Array();
        for(var data in this.props.info ) {
            var tmp = {};
            tmp['key'] = data;
            tmp['value'] = this.props.info[data];
            infoUserTmp.push(tmp);
        }
        // Map the Array to create a special input for each one
        var NodeInput = infoUserTmp.map(function(info) {
            return (
                <input type="text" value={info.value} 
                    id={info.key}
                    onChange={this.handleChangeForm}
                />
            );
        }, this);
        return (
            <div className="userForm" >
                <form onSubmit={this.handleSubmit} >
                    {NodeInput}
                    <br />
                    <input type="password"
                    />
                    <input type="password"
                    />
                    <button>Okay</button>
                </form>
            </div>
        );
    }
});

module.exports = UserForm;
