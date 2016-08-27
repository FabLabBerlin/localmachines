var $ = require('jquery');
var reactor = require('../reactor');
var getters = require('../getters');
var GlobalActions = require('../actions/GlobalActions');
var HeaderNav = require('./Header/HeaderNav');
var LoaderLocal = require('./LoaderLocal');
var LocationActions = require('../actions/LocationActions');
var LoginActions = require('../actions/LoginActions');
var LoginStore = require('../stores/LoginStore');
var MachinePage = require('./MachinePage/MachinePage');
var React = require('react');
var toastr = require('../toastr');

import {hashHistory} from 'react-router';

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';


/*
 * App
 * @component:
 * manage the router
 * it is in every page
 * navigation bar and footer are in this component
 */
 var App = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isLoading: getters.getIsLoading,
      isLogged: getters.getIsLogged
    };
  },

  componentWillMount() {
    const isLogged = reactor.evaluateToJS(getters.getIsLogged);

    if (!isLogged && this.props.location.pathname !== '/product') {
      LoginActions.tryAutoLogin(this.context.router, {
        loggedIn: () => {
          if (this.props.location.pathname === '/login') {
            hashHistory.push('/machine');
          }
        },
        loggedOut: () => {
          if (this.props.location.pathname !== '/login') {
            hashHistory.push('/login');
          }
        }
      });
      LocationActions.loadLocations();
    } else {
      if (this.props.location.pathname === '/login') {
        hashHistory.push('/machines');
      }
    }
  },

  render() {
    const footerAbsoluteBottom = !this.state.isLogged &&
      this.props.location.pathname !== '/product';

    // TODO: wait for ajax promise or something similar
    //       ... until then this below is just a quick fix:
    if (!this.state.isLogged && !(
        this.props.location.pathname === '/login' ||
        this.props.location.pathname === '/product')) {
      window.location.href = '/machines/#/login';
      return <LoaderLocal/>;
    }

    return (
      <div className="app">
        {this.props.location.pathname !== '/product' ? <HeaderNav location={this.props.location}/> : null}
        <div id="main-content">
          {(this.state.isLogged ||
            this.props.location.pathname === '/login' ||
            this.props.location.pathname === '/product' ||
            this.props.location.pathname.indexOf('forgot_password') > 0
            ) ? 
            (this.props.children || <MachinePage/>) :
            <LoaderLocal />
          }
        </div>
        {this.props.location.pathname !== '/product' ? (
          <footer className={footerAbsoluteBottom ? 'absolute-bottom' : ''}>
            <div className="container-fluid">
              <div className="col-md-4 text-center">
                <i className="fa fa-copyright"></i> Makea Industries GmbH 2016
              </div>
              <div className="col-md-4 text-center">
                In case you are interested in using EASY LAB in your own
                Lab, <a href="/machines/#/product">visit the Product Page</a>.
              </div>
              <div className="col-md-4 text-center">
                <a href="https://fablab.berlin/en/content/2-Imprint">Imprint</a>
              </div>
            </div>
          </footer>
        ) : null}
        {
          this.state.isLoading ?
          (
            <div id="loader-global">
              <div className="spinner">
                <i className="fa fa-cog fa-spin"></i>
              </div>
            </div>
          )
          : ''
        }
      </div>
    );

  },

  signupNewsletter() {
    VexDialog.prompt({
      message: 'Please enter your E-Mail address:',
      placeholder: 'E-Mail',
      callback: (value) => {
        if (value) {
          var email = value;

          GlobalActions.performSubscribeNewsletter(email);
        } else if (value !== false) {
          toastr.error('No token');
        }
      }
    });
  }

});

export default App;
