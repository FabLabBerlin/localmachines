var $ = require('jquery');

var getters = require('../../getters');
var Location = require('../../modules/Location');
var LoginStore = require('../../stores/LoginStore');

var LoginActions = require('../../actions/LoginActions');

var React = require('react');
var ReactDOM = require('react-dom');
var reactor = require('../../reactor');

var _ = require('lodash');

import {hashHistory} from 'react-router';


/*
 * Login component
 * Handle the login page
 * Give permission to edit your information
 */
var Login = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: Location.getters.getLocation,
      locationId: Location.getters.getLocationId,
      locations: Location.getters.getLocations
    };
  },

  /*
   * Submit the form
   * Clear the input
   */
  handleSubmit(event) {
    event.preventDefault();
    var data = {
      username: this.refs.name.value,
      password: this.refs.password.value,
      location: parseInt(this.refs.location.value)
    };
    LoginActions.submitLoginForm(data, this.context.router);
  },

  /*
   * If you are already connected, will skip the page
   * listen to the onChange event from the UserStore
   */
  componentDidMount() {
    LoginActions.tryAutoLogin(this.context.router, {});

    if (reactor.evaluateToJS(getters.getIsLogged)) {
      hashHistory.push('/machine');
    }
  },

  /*
   * Render the form and the button inside of the App component
   */
  render() {
    var locations;
    if (this.state) {
      try {
        locations = this.state.locations;
      } catch (e) {
        locations = [];
      }
    } else {
      locations = [];
    }

    if (locations) {
    return (
      <form className="login-form" method="post" onSubmit={this.handleSubmit}>
        <div className="regular-login">

          <h2 className="login-heading">Please log in</h2>
          <input
            ref="name"
            type="text"
            name="username"
            className="form-control"
            placeholder="E-Mail"
            required
            autoFocus="on"
            autoCorrect="off"
            autoCapitalize="off"
          />
          <input
            type="password"
            name="password"
            className="form-control"
            ref="password"
            placeholder="Password"
            required
          />
          <select
            className="form-control location-picker"
            ref="location"
            name="location"
            value={this.state.locationId}
            onChange={this.updateLocation}
            required>
            {locations.map((location, i) => {
              if (location.get('Approved')) {
                return (
                  <option key={i} value={location.get('Id')}>
                    {location.get('Title')}
                  </option>
                );
              }
            })}
          </select>
          <button className="btn btn-primary btn-block btn-login"
            type="submit">Log In</button>

          {this.state.location ? (
            <div className="signup-link">
              Do not have an account yet? <a href={'/signup/#/form?location=' + this.state.location.get('Id')}
                onClick={this.goToSignUp}>
                Sign up for {this.state.location.get('Title')}
              </a> now!
            </div>
          ) : null}

          <div className="text-center">
            Forgot your password? <a href="/machines/#/forgot_password/start">
              Oops...
            </a>
          </div>

        </div>
      </form>
    );
    } else {
      return (<div></div>);
    }
  },

  updateLocation() {
    var locationId = parseInt(this.refs.location.value);
    Location.actions.setLocationId(locationId);
  }
});

export default Login;
