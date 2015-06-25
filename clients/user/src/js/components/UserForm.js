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
                <div className="col-md-6" >
                    <div className="form-group" >
                        <label htmlFor="user-information" >{info.key}</label>
                        <input type="text" value={info.value} 
                            id={info.key}
                            className="form-control"
                            onChange={this.handleChangeForm}
                        />
                    </div>
                </div>
            );
        }, this);
        return (
            //l61, put font-awesome to get the logo
            <div className="" >
                <form onSubmit={this.handleSubmit} >
                    <div className="row" >
                        {NodeInput}
                    </div>
                    <div className="row" >
                        <div className="form-group" >
                            <div className="col-sm-6" >
                                <label htmlFor="user-password" >User Password </label>
                                <input 
                                    type="password" className="form-control"
                                    placeholder="new password"
                                />
                            </div>
                            <div className="col-sm-6" >
                                <label htmlFor="user-password" >User Password </label>
                                <input 
                                    type="password" className="form-control"
                                    placeholder="repeat password"
                                />
                            </div>
                        </div>
                    </div>
                    <div className="col-sm-6" >
                        <button className="btn btn-primary btn-lg">
                            Save
                        </button>
                    </div>
                </form>
            </div>
        );
    }
});

module.exports = UserForm;
