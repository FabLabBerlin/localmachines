import React from 'react';
import {Link, RouteHandler} from 'react-router';

var App = React.createClass({
    render: function() {
        return (
            <div className="app" >
                <header>
                    <div className="container-fuild">
                        <img src="assets/logo_fablab_berlin.svg" className="brand-image" />
                    </div>
                </header>

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
