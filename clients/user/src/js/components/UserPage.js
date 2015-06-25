import React from 'react';
import {Navigation} from 'react-router';
import UserForm from './UserForm';
import MachineList from './MachineList';
import Membership from './Membership';
import UserActions from '../actions/UserActions';
import UserStore from '../stores/UserStore'

var UserPage = React.createClass({

    mixins: [ Navigation ],

    // If not login, redirect to the login page
    statics: {
        willTransitionTo(transition) {
            if(!UserStore.getIsLogged()) {
                transition.redirect('login');
            }
        }
    },

    /*
     * Getting stuff from OUTSIDE of the component
     */
    // getting state from UserStore
    getInitialState: function() {
        var _infoUser = UserStore.getInfoUser();
        var _infoMachine = UserStore.getInfoMachine();
        var _infoMembership = UserStore.getMembership();
        return {
            infoUser: _infoUser,
            infoMachine: _infoMachine,
            infoMembership: _infoMembership
        };
    },

    // Pass the responsabilité to the store via the action
    handleSubmit() {
        UserActions.submitState(this.state.infoUser);
    },

    handleLogout(){
        UserActions.logout();
    },

    onChangeLogout() {
        if( !UserStore.getIsLogged() ){
            this.replaceWith('login');
        }
    },

    /*
     * INSIDE the component
     */
    // Change the state of the input related
    handleChangeForm(event) {
        // Create a temporary state to replace the old one
        var tmpState = this.state.infoUser;
        tmpState[event.target.id] = event.target.value;
        this.setState({
            infoUser: tmpState
        });
    },

    componentDidMount() {
        UserStore.onChangeLogout = this.onChangeLogout;
    },

    render() {
        return (
            <div className="signup" >
                <div className="container-fluid" >
                    <div className="signup-form" >
                        <button onClick={this.handleLogout} >Logout</button>
                        <UserForm info={this.state.infoUser} 
                            func={this.handleChangeForm}
                            submit={this.handleSubmit}
                        />
                        <MachineList info={this.state.infoMachine} />
                        <Membership info={this.state.infoMembership} />
                    </div>
                </div>
            </div>
        );
    }
});

module.exports = UserPage;
