var reactor = require('../reactor');
var getters = require('../getters');
var HeaderNav = require('./Header/HeaderNav');
var LoaderLocal = require('./LoaderLocal');
var LocationActions = require('../actions/LocationActions');
var LoginActions = require('../actions/LoginActions');
var LoginStore = require('../stores/LoginStore');
var {Navigation} = require('react-router');
var React = require('react');
var RouteHandler = require('react-router').RouteHandler;


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
    if (!isLogged) {
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
        <HeaderNav />
        {(this.state.isLogged || window.location.hash === '#/login') ? 
          <RouteHandler /> :
          <LoaderLocal />
        }
        <footer className={this.state.isLogged ? '' : 'absolute-bottom'}>
          <div className="container-fluid">
            <div className="col-md-4 text-center">
              <i className="fa fa-copyright"></i> Makea Industries GmbH 2016
            </div>
            <div className="col-md-4 text-center"></div>
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

  }
});

export default App;
