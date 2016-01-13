var Global = require('../modules/Global');
var HeaderNav = require('./HeaderNav');
var Login = require('../modules/Login');
var React = require('react');
var reactor = require('../reactor');
var RouteHandler = require('react-router').RouteHandler;


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
    console.log('App.js: getDataBindings');
    console.log('Login.getters:', Login.getters);
    console.log('Login.getters.getIsLoading:', Login.getters.getIsLoading);
    return {
      isLoading: Global.getters.getIsLoading
    };
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
    console.log('App.js: render');
    
    const isLogged = reactor.evaluateToJS(Login.getters.getIsLogged);
    
    return (
      <div className="app">
        <HeaderNav />
        <RouteHandler />
        <footer className={isLogged ? '' : 'absolute-bottom'}>
          <div className="container-fluid">
            <i className="fa fa-copyright"></i> Fab Lab Berlin 2015
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
