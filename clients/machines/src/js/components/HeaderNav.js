var LoginActions = require('../actions/LoginActions');
var {Navigation} = require('react-router');
var React = require('react');


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


var HeaderNav = React.createClass({
  mixins: [ Navigation ],

  handleClick() {
    LoginActions.logout(this.context.router);
  },

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
      <div>

<nav className="navbar navbar-default">
  <div className="container-fluid">
    
    <div className="navbar-header">
      <button type="button" 
        className="navbar-toggle collapsed" 
        data-toggle="collapse" 
        data-target="#machines-navbar" 
        aria-expanded="false">
        
        <span className="sr-only">Toggle navigation</span>
        <span className="icon-bar"></span>
        <span className="icon-bar"></span>
        <span className="icon-bar"></span>
      </button>
      
      <div className="navbar-brand">
        <img src="img/logo_easylab.svg" className="brand-image"/>
        <img src="img/logo_small.svg" className="brand-image brand-image-mobile"/>
      </div>
    </div>

    <div className="collapse navbar-collapse" id="machines-navbar">
      
      <ul className="nav navbar-nav navbar-right">
        
        {buttons}
        
        <li>
          <a href="#"
            onClick={this.handleClick}>
            <i className="fa fa-sign-out"></i> Sign Out
          </a>
        </li>

      </ul>
      
    </div>
  </div>
</nav>

</div>

    );
  }
});

export default HeaderNav;
