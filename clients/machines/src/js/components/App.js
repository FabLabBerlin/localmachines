var $ = require('jquery');
var reactor = require('../reactor');
var getters = require('../getters');
var GlobalActions = require('../actions/GlobalActions');
var HeaderNav = require('./Header/HeaderNav');
var LoaderLocal = require('./LoaderLocal');
var LocationActions = require('../actions/LocationActions');
var LoginActions = require('../actions/LoginActions');
var LoginStore = require('../stores/LoginStore');
var {Navigation} = require('react-router');
var React = require('react');
var RouteHandler = require('react-router').RouteHandler;
var toastr = require('../toastr');

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

  mixins: [ Navigation, reactor.ReactMixin ],

  getDataBindings() {
    return {
      isLoading: getters.getIsLoading,
      isLogged: getters.getIsLogged
    };
  },

  componentWillMount() {
    const isLogged = reactor.evaluateToJS(getters.getIsLogged);
    if (!isLogged && window.location.hash !== '#/product') {
      LoginActions.tryPassLoginForm(this.context.router, {
        loggedIn: () => {
          if (window.location.hash === '#/login') {
            this.transitionTo('/machine');
          }
        },
        loggedOut: () => {
          if (window.location.hash !== '#/login') {
            this.transitionTo('/login');
          }
        }
      });
      LocationActions.loadLocations();
    } else {
      if (window.location.hash === '#/login') {
        this.transitionTo('/machines');
      }
    }
  },

  /*
   * Render:
   *  - navBar
   *  - all the component which are under the router control
   *  - footer
   * If user is logged, display a exit button
   * If he's logged and there is no nfc port, can switch to user interface
   */
  render() {
    return (
      <div className="app">
        {window.location.hash !== '#/product' ? <HeaderNav /> : null}
        <div id="main-content">
          {(this.state.isLogged ||
            window.location.hash === '#/login' ||
            window.location.hash === '#/product' ||
            window.location.hash.indexOf('forgot_password') > 0
            ) ? 
            <RouteHandler /> :
            <LoaderLocal />
          }
        </div>
        <footer className={this.state.isLogged ? '' : 'absolute-bottom'}>
          <div className="container-fluid">
            <div className="col-md-4 text-center">
              <i className="fa fa-copyright"></i> Makea Industries GmbH 2016
            </div>
            <div className="col-md-4 text-center">
              In case you are interested in using EASY LAB in your own
              Lab, <a href="javascript:void(0);" onClick={this.signupNewsletter}>signup to our newsletter</a>.
            </div>
            <div className="col-md-4 text-center">
              <a href="https://fablab.berlin/en/content/2-Imprint">Imprint</a>
            </div>
          </div>
        </footer>
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

          this.performSubscribe(email);
        } else if (value !== false) {
          toastr.error('No token');
        }
      }
    });
  },

  performSubscribe(email) {
    GlobalActions.showGlobalLoader();
    $.ajax({
      method: 'POST',
      url: '/api/newsletters/easylab_dev',
      data: {
        email: email
      },
      success: function() {
        GlobalActions.hideGlobalLoader();
        toastr.info('Please check your E-Mails to confirm the subscription.');
      },
      error: function() {
        GlobalActions.hideGlobalLoader();
        toastr.error('An error occurred, please try again later.');
      }
    });
  }

});

export default App;
