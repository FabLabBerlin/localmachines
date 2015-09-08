var LoginActions = require('../actions/LoginActions');
var {Navigation} = require('react-router');
var React = require('react');
var Reactor = require('../reactor');
var Getters = require('../getters');


var MenuItem = React.createClass({
  render() {
    
    var activeClass = '';
    if (this.props.href === (window.location.pathname + window.location.hash)) {
      activeClass = 'active';
    }

    return (
      <li className={activeClass}>
        <a id={this.props.id} href={this.props.href}>
          <i className={'fa ' + 
            this.props.faIconClass}></i> {this.props.label}
        </a>
      </li>
    );
  }
});

var BurgerMenuToggle = React.createClass({
  render() {
    return (
      <button type="button" 
        className="navbar-toggle collapsed" 
        data-toggle="collapse" 
        data-target="#burger-menu" 
        aria-expanded="false">
        
        <span className="sr-only">Burger Menu</span>
        <span className="icon-bar"></span>
        <span className="icon-bar"></span>
        <span className="icon-bar"></span>
      </button>
    );
  }
});

var MainMenu = React.createClass({
  render() {
    var buttons = [];

    if (!window.libnfc) {
      buttons.push(
        <MenuItem href="/machines/#/machine"
          faIconClass="fa-wrench"
          label="Machines">
        </MenuItem>);
  
      buttons.push(
        <MenuItem href="/machines/#/profile"
          faIconClass="fa-user"
          label="Profile">
        </MenuItem>);
      
      buttons.push(
        <MenuItem href="/machines/#/spendings"
          faIconClass="fa-money"
          label="Spendings">
        </MenuItem>);
    }

    return (
      <div className="collapse navbar-collapse" id="burger-menu">
        <ul className="nav navbar-nav navbar-right">
          {buttons}
          <li>
            <a href="#" onClick={this.handleClick} className="sign-out">
              <i className="fa fa-sign-out"></i> Sign Out
            </a>
          </li>
        </ul>  
      </div>
    );
  }
});

var HeaderNav = React.createClass({
  
  mixins: [ Navigation ],

  handleClick() {
    LoginActions.logout(this.context.router);
  },

  render() {
    const isLogged = Reactor.evaluateToJS(Getters.getIsLogged);

    return (
      <div>

        <nav className="navbar navbar-default">
          <div className="container-fluid">
            
            <div className="navbar-header">
              
              {isLogged ? (<BurgerMenuToggle />) : ('')}
              
              <div className="navbar-brand">
                <img src="img/logo-easylab.svg" className="brand-image hidden-xs"/>
                <img src="img/logo-small.svg" className="brand-image visible-xs-block"/>
              </div>

            </div>
        
            {isLogged ? (<MainMenu />) : ('')}

          </div>
        </nav>

      </div>
    );
  }

});

export default HeaderNav;
