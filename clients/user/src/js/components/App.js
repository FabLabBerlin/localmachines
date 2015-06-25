import UserStore from '../stores/UserStore';
import UserActions from '../actions/UserActions';
import React from 'react';
import {Link, RouteHandler} from 'react-router';

var App = React.createClass({
    render: function() {
        return (
            <div className="app" >
                <nav class="navbar navbar-default" >
                    <div className="collapse navbar-collapse" >
                        <ul className="nav navbar-nav" >
                            <li>
                                <img src="assets/logo_fablab_berlin.svg" className="brand-image" />
                            </li>
                        </ul>
                        {UserStore.getIsLogged() ? (
                            <ul className="nav navbar-nav navbar-right" >
                                <li>
                                    <a onClick={UserActions.logout} >
                                        <i className="glyphicon glyphicon-log-out"></i>
                                    </a>
                                </li>
                            </ul>
                            ):(
                            ''
                            )}
                        </div>
                    </nav>

                    <RouteHandler />

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
