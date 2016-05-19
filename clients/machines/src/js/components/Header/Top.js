var getters = require('../../getters');
var LocationGetters = require('../../modules/Location/getters');
var React = require('react');
var reactor = require('../../reactor');


var UserRole = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      userLocation: LocationGetters.getUserLocation
    };
  },

  render() {
    return (
      <div className="nav-user-role">
        {this.userRole()}
      </div>
    );
  },

  userRole() {
    if (this.state.userLocation) {
      switch (this.state.userLocation.get('UserRole')) {
      case 'admin':
        return 'Admin';
      case 'staff':
        return 'Staff';
      default:
        return 'Member';
      }
    }
  }

});


var Right = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      user: getters.getUser,
      userLocation: LocationGetters.getUserLocation
    };
  },

  render() {
    return (
      <div className="nav-right">
        {this.state.user ? (
          <div>
            <div className="nav-user-name">
              {this.state.user.get('FirstName')} {this.state.user.get('LastName')}
            </div>
            <UserRole/>
          </div>
        ) : null}
      </div>
    );
  }

});


var Top = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      user: getters.getUser
    };
  },

  render() {
    return (
      <div className="row">
        <div className="col-xs-12">
          <div className="nav-top">
            <a className="nav-logo" 
               href="/machines/#/machine">
              <img src="/machines/assets/img/logo-easylab.svg" className="brand-image hidden-xs"/>
              <img src="/machines/assets/img/logo-small.svg" className="brand-image visible-xs-block"/>
            </a>
            <div className="nav-title">
              {this.state.location ? (
                <span className="nav-title-lab">
                  {this.state.location.get('Title')}
                </span>
              ) : null}
              <span className="nav-title-easylab">EASY LAB</span>
            </div>
            <Right/>
          </div>
        </div>
      </div>
    );
  }
});

export default Top;
