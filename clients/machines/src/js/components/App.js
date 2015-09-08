var reactor = require('../reactor');
var getters = require('../getters');
var HeaderNav = require('./HeaderNav');
var LoginStore = require('../stores/LoginStore');
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
  /*
   * Render:
   *  - navBar
   *  - all the component which are under the router control
   *  - footer
   * If user is logged, display a exit button
   * If he's logged and there is no nfc port, can switch to user interface
   */
  render() {
    
    const isLogged = reactor.evaluateToJS(getters.getIsLogged);
    
    return (
      <div className="app">
        <HeaderNav />
        <RouteHandler />
        <footer className={isLogged ? '' : 'absolute-bottom'}>
          <div className="container-fluid">
            <i className="fa fa-copyright"></i> Fab Lab Berlin 2015
          </div>
        </footer>
      </div>
    );

  }
});

export default App;
