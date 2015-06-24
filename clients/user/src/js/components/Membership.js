import React from 'react';

var Membership = React.createClass({
    render() {
        var membership = this.props.info[0];
        return (
            <div className="membership" >
                <p> Membership </p>
                <ul>
                    <li>Membership Id : {membership.MembershipId}</li>
                    <li>Start date : {membership.StartDate}</li>
                </ul>
            </div>
        );
    }
});

module.exports = Membership;
