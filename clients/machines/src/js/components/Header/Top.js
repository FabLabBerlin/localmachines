var $ = require('jquery');
var getters = require('../../getters');
var GlobalActions = require('../../actions/GlobalActions');
var LocationGetters = require('../../modules/Location/getters');
var LoginActions = require('../../actions/LoginActions');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('toastr');


var Right = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isAdmin: LocationGetters.getIsAdmin,
      user: getters.getUser,
      userLocation: LocationGetters.getUserLocation
    };
  },

  handleLogout() {
    toastr.info('Bye! Don\'t forget to turn off the machines when done!');
    GlobalActions.showGlobalLoader();
    window.setTimeout(() => {
      window.location.href = '/logout';
    }, 4000);
  },

  render() {
    if (this.state.user) {
      return (
        <div className="nav-right pull-right">
          <div className="nav-dropdown">
            <div className="dropdown">
              <button className="btn btn-default dropdown-toggle" type="button" id="dropdownMenu1" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <div className="nav-user">
                  <div className="nav-user-name">
                    {(this.state.user.get('FirstName') || '') + ' ' + (this.state.user.get('LastName') || '')}
                  </div>
                  <div className="nav-user-role">
                    {this.userRole()}
                  </div>
                </div>
                <div id="nav-caret-container">
                  <span className="caret"></span>
                </div>
              </button>
              <ul className="dropdown-menu pull-right" aria-labelledby="dropdownMenu1">
                {this.state.isAdmin ? (
                  <li>
                    <a href="/machines/#/admin/categories">Categories</a>
                  </li>
                ) : null}
                <li><a href="/machines/#/profile">Info</a></li>
                {this.state.isAdmin ? (
                  <li>
                    <a href="/machines/#/admin/invoices">Invoices</a>
                  </li>
                ) : null}
                {this.state.user.get('SuperAdmin') ? (
                  <li>
                    <a href="/machines/#/admin/locations">Locations</a>
                  </li>
                ) : null}
                {this.state.isAdmin ? (
                  <li>
                    <a href="/machines/#/admin/settings">Settings</a>
                  </li>
                ) : null}
                <li><a href="/machines/#/feedback">Feedback</a></li>
                <li role="separator" className="divider"></li>
                <li><a onClick={this.handleLogout} href="javscript:void(0);">Log out</a></li>
              </ul>
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
      locationId: LocationGetters.getLocationId,
      user: getters.getUser
    };
  },

  render() {
    return (
      <div className="nav-top row">
        <div className="nav-left" style={{overflow: 'hidden'}}>
          <a className="nav-logo" 
             href="/machines/#/machine">
            {this.state.location ? (
              this.state.location.get('Logo') ?
                <img src={('/files/' + this.state.location.get('Logo'))}/> :
                <img src="/machines/assets/img/logo-easylab.svg"/>
            ) : null}
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
        <Right/>
      </div>
    );
  }
});

export default Top;
