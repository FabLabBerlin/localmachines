import React from 'react';
import {Link, RouteHandler} from 'react-router';

var App = React.createClass({
    render: function() {
        return (
            <div className="app" >
                <RouteHandler />
            </div>
        );
    }
});

module.exports = App;
