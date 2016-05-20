var getters = require('../../getters');
var LocationGetters = require('../../modules/Location/getters');
var React = require('react');
var reactor = require('../../reactor');


var Right = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      user: getters.getUser,
      userLocation: LocationGetters.getUserLocation
    };
  },

  render() {
    if (this.state.user) {
      return (
        <div className="nav-right pull-right">
          <div className="nav-dropdown">
            <i className="fa fa-caret-down"/>
          </div>
          <div className="nav-user">
            <div className="nav-user-name">
              {this.state.user.get('FirstName')} {this.state.user.get('LastName')}
            </div>
            <div className="nav-user-role">
              {this.userRole()}
            </div>
          </div>
        </div>
      );
    } else {
      return <div className="nav-right"/>;
    }
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
      <div className="nav-top row">
        <div className="col-xs-6 nav-left">
          <a className="nav-logo" 
             href="/machines/#/machine">
            <img src="/machines/assets/img/logo-easylab.svg"/>
          </a>
          <div className="nav-title hidden-xs">
            {this.state.location ? (
              <span className="nav-title-lab">
                {this.state.location.get('Title')}
              </span>
            ) : null}
            <span className="nav-title-easylab">EASY LAB</span>
          </div>
        </div>
        <div className="col-xs-6 pull-right">
          <Right/>
        </div>
      </div>
    );
  }
});

export default Top;
