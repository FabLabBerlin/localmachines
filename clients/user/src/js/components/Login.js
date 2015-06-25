import React from 'react';
import {Navigation} from 'react-router';

import UserStore from '../stores/UserStore';
import LoginActions from '../actions/LoginActions';

var Login = React.createClass({

    mixins: [ Navigation ],

    /*
     * Getting stuff from OUTSIDE of the component
     */

    // If not login, redirect to the login page
    statics: {
        willTransitionTo(transition) {
            if(UserStore.getIsLogged()) {
                transition.redirect('login');
            }
        }
    },

    componentDidMount() {
        LoginActions.submitLoginForm(this.state);
        UserStore.onChange = this.onChange;
    },

    // Sending the form
    handleSubmit: function(event) {
        event.preventDefault();
        LoginActions.submitLoginForm(this.state);
        this.clearAndFocus();
    },

    onChange() {
        if( UserStore.getIsLogged ) {
            this.replaceWith('user');
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
            <div class="login" >
                <div className="container-fluid">
                    <div className="regular-login" >
                        <form className="login-form"
                            onSubmit={this.handleSubmit} >
                            <h2 class="login-heading">Please log in</h2>
                            <input 
                                ref="name" 
                                type="text"
                                className="form-control"
                                value={this.state.username}
                                onChange={this.handleChange}
                                placeholder="Username" 
                                required
                                autofocus
                            />
                            <input
                                type="password" 
                                className="form-control"
                                ref="Password" 
                                value={this.state.password}
                                onChange={this.handleChange}
                                placeholder="password" 
                                required
                            />
                            <button className="btn btn-primary btn-block btn-login"
                                type="submit">Log In</button>
                        </form>
                    </div>
                    </div>
                    </div>
        );
    }
});

module.exports = Login;
