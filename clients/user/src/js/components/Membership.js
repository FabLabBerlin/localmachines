import React from 'react';

var Membership = React.createClass({
    render() {
        return (
            <div className="membership" >
                <p> Membership </p>
                <ul>
                    <li>here the membership</li><!-- here change to be dynamic -->
                </ul>
            </div>
        );
    }
});

module.exports = Membership;
