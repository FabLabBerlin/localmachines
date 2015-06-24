import React from 'react';

var Membership = React.createClass({
    render() {
        if(this.props.info.length != 0) {
            var membership = this.props.info[0];
            var MembershipNode = 
                <ul>
                    <li> Membership id: {membership.MembershipId}</li>
                    <li>Start date: {membership.StartDate}</li>
                </ul>
        } else {
            var MembershipNode = <p>You do not have any membership</p>
        }
        return (
            <div className="membership" >
                <p> Membership </p>
                {MembershipNode}
            </div>
        );
    }
});

module.exports = Membership;
