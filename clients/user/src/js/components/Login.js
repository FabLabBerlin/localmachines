import React from 'react';
import {Navigation} from 'react-router';

import UserStore from '../stores/UserStore';
import LoginActions from '../actions/LoginActions';

var Login = React.createClass({

    mixins: [ Navigation ],

    /*
     * Getting stuff from OUTSIDE of the component
     */

    // Sending the form
    handleSubmit: function(event) {
        event.preventDefault();
        LoginActions.submitLoginForm(this.state);
        this.clearAndFocus();
        if( UserStore.getIsLogged ) {
            this.transitionTo('/');
        }
    },

    /*
     * INSIDE the component
     */
    // Change the state of the input related
    getInitialState: function() {
        return {
            username: '',
            password: ''
        };
    },

    // To save Username and password before sending them
    handleChange: function(e) {
        this.setState({
            [e.target.name]: e.target.value
        });
    },
    // Clear the state before sending the form
    clearAndFocus: function() {
        this.setState({username: '', password: ''}, function() {
            React.findDOMNode(this.refs.name).focus();
        });
    },

    render() {
        return (
            <div className="login" >
                <p>Logging page</p>
                <form onSubmit={this.handleSubmit} >
                    <input 
                        ref="name" 
                        name="username"
                        value={this.state.username}
                        onChange={this.handleChange}
                        placeholder="username" 
                    />
                    <input
                        type="password" 
                        name="password"
                        ref="password" 
                        value={this.state.password}
                        onChange={this.handleChange}
                        placeholder="password" 
                    />
                    <button > Submit </button>
                </form>
            </div>
        );
    }
});

module.exports = Login;
