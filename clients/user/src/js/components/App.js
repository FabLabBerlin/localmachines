import UserStore from '../stores/UserStore';
import UserActions from '../actions/UserActions';
import React from 'react';
import {Link, RouteHandler} from 'react-router';

/*
 * App
 * @component: 
 * manage the router
 * it is in every page
 * navigation bar and footer are in this component
 */
var App = React.createClass({
  /*
   * Render:
   *  - navBar
   *  - all the component which are under the router control
   *  - footer
   */
  render: function() {
    return (
      <div className="app">

        <header>
          <div className="container-fluid">

            <div className="row">
              <div className="col-xs-6">
                <img src="img/logo_fablab_berlin.svg" 
                     className="brand-image"/>
              </div>
              {UserStore.getIsLogged() ? (
                <div className="col-xs-6 text-right">
                  <button 
                    className="btn btn-danger btn-logout pull-right" 
                    onClick={UserActions.logout}>
                    <i className="fa fa-sign-out"></i>
                  </button>
                  <a href="/machines"
                    className="btn btn-info linkToPanel"
                    role="button" >
                    Switch to <br/>
                    machine panel
                  </a>
                </div>
              ):('')}
            </div>

          </div>
        </header>

        <RouteHandler />

        <footer>
          <div className="container-fluid">
            <i className="fa fa-copyright"></i> Fab Lab Berlin 2015
          </div>
        </footer>

      </div>
    );
  }
});

module.exports = App;
