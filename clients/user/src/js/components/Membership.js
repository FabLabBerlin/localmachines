import React from 'react';

var Membership = React.createClass({
    render() {
        if(this.props.info.length != 0) {
            var membership = this.props.info[0];
            var MembershipNode = 
                <tr>
                    <td> Membership id: {membership.MembershipId}</td>
                    <td>Start date: {membership.StartDate}</td>
                </tr>
        } else {
            var MembershipNode = <tr>You do not have any membership</tr>
        }
        return (
            <table className="table table-striped table-hover" >
                <thead>
                    <tr>
                        <th>Membership Id</th>
                        <th>Start date</th>
                    </tr>
                </thead>
                <tbody>
                    {MembershipNode}
                </tbody>
            </table>
        );
    }
});

module.exports = Membership;
