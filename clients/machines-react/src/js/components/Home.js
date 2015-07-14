import React from 'react';
import UserStore from '../stores/UserStore';

var Home = React.createClass({

    statics: {
        willTransitionTo(transition) {
            if(!UserStore.getIsLogged()) {
                transition.redirect('login');
            }
        }
    },

    render: function() {
        return (
            <div>Home</div>
        );
    }
});

module.exports = Home;
