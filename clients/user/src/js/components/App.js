import React from 'react';
import {Link, RouteHandler} from 'react-router';

var App = React.createClass({
    render: function() {
        return (
            <div className="app" >
                <header>
                    <div class="container-fuild">
                        <img src="assets/logo_fablab_berlin.svg" className="brand-image" />
                    </div>
                </header>

                <RouteHandler />

                <footer class="absolute-bottom" >
                    <div class="container-fuild">
                        <i class="fa fa-copyright"></i> Fab Lab Berlin 2015
                    </div>
                </footer>
            </div>
        );
    }
});

module.exports = App;
