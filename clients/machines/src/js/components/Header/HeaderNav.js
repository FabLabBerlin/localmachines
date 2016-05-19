var Bottom = require('./Bottom');
var getters = require('../../getters');
var LoginActions = require('../../actions/LoginActions');
var {Navigation} = require('react-router');
var React = require('react');
var reactor = require('../../reactor');
var Top = require('./Top');


var HeaderNavBrand = React.createClass({
  render() {
    return (
      <div className="navbar-brand">
        <a 
          className="brand-link" 
          href="/machines/#/machine">
          <img src="/machines/assets/img/logo-easylab.svg" className="brand-image hidden-xs"/>
          <img src="/machines/assets/img/logo-small.svg" className="brand-image visible-xs-block"/>
        </a>
      </div>
    );
  }
});

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

  mixins: [ Navigation, reactor.ReactMixin ],

  getDataBindings() {
    return {
      user: getters.getUser
    };
  },

  signOut() {
    LoginActions.logout(this.context.router);
  },

  render() {
    var buttons = [];

    if (!window.libnfc) {
      buttons.push(
        <MenuItem key={1}
          href="/machines/#/machine"
          faIconClass="fa-plug"
          label="Machines">
        </MenuItem>,
        <MenuItem key={2}
          href="/machines/#/reservations"
          faIconClass="fa-calendar-check-o"
          label="Reservations">
        </MenuItem>,      
        <MenuItem key={3}
          href="/machines/#/spendings"
          faIconClass="fa-money"
          label="Spendings">
        </MenuItem>,
        <MenuItem key={4}
          href="/machines/#/feedback"
          faIconClass="fa-paper-plane"
          label="Feedback">
        </MenuItem>,
        <MenuItem key={5}
          href="/machines/#/profile"
          faIconClass="fa-user"
          label={this.state.user.get('FirstName')}>
        </MenuItem>);
    }

    return (
      <div className="collapse navbar-collapse" id="burger-menu">
        <ul className="nav navbar-nav navbar-right">
          {buttons}
          <li>
            <a href="#" onClick={this.signOut} className="sign-out">
              <i className="fa fa-sign-out"></i> Sign out
            </a>
          </li>
        </ul>  
      </div>
    );
  }
});


var OldHeaderNav = React.createClass({
  render() {
    const isLogged = reactor.evaluateToJS(getters.getIsLogged);

    return (
      <div>
        <nav className="navbar navbar-default">
          <div className="container-fluid">
            <div className="navbar-header">              
              {isLogged ? (<BurgerMenuToggle />) : ('')}   
              <HeaderNavBrand />
            </div>
            {isLogged ? (<MainMenu />) : ('')}
          </div>
        </nav>
      </div>
    );
  }
});


var HeaderNav = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isLogged: getters.getIsLogged
    }
  },

  render() {
    if (this.state.isLogged) {
      return (
        <div>
          <Top/>
          <Bottom/>
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default HeaderNav;
