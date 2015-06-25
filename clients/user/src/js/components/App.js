import UserStore from '../stores/UserStore';
import UserActions from '../actions/UserActions';
import React from 'react';
import {Link, RouteHandler} from 'react-router';

var App = React.createClass({
    render: function() {
        return (
            <div className="app" >
                <nav className="navbar navbar-default" >
                    <div className="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
                        <ul className="nav navbar-nav" >
                            <li>
                                <img src="assets/logo_fablab_berlin.svg" className="brand-image" />
                            </li>
                        </ul>
                        {UserStore.getIsLogged() ? (
                            <ul className="nav navbar-nav navbar-right" >
                                <li>
                                    <button 
                                        onClick={UserActions.logout}
                                        className="btn btn-danger btn-lg2"
                                        >
                                        Logout
                                    </button>
                                </li>
                            </ul>
                            ):('')}
                        </div>
                    </nav>

                    <div className="blue" >
                        <RouteHandler />
                    </div>

                    <footer className="absolute-bottom" >
                        <div className="container-fuild">
                            <i className="fa fa-copyright"></i> Fab Lab Berlin 2015
                        </div>
                    </footer>
                </div>
        );
    }
});

module.exports = App;
