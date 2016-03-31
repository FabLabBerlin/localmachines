var $ = require('jquery');

var getters = require('../../getters');

var LoginStore = require('../../stores/LoginStore');
var LocationStore = require('../../stores/LocationStore');

var LoginActions = require('../../actions/LoginActions');
var LocationActions = require('../../actions/LocationActions');

var {Navigation} = require('react-router');

var React = require('react');
var reactor = require('../../reactor');

var _ = require('lodash');


/*
 * Login component
 * Handle the login page
 * Give permission to edit your information
 */
var Login = React.createClass({

  /*
   * To use transitionTo/replaceWith/redirect and some function related to the router
   */
  mixins: [ Navigation, reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: getters.getLocation,
      locations: getters.getLocations
    };
  },

  /*
   * Submit the form
   * Clear the input
   */
  handleSubmit(event) {
    event.preventDefault();
    var data = {
      username: this.refs.name.getDOMNode().value,
      password: this.refs.password.getDOMNode().value,
      location: this.state.location.Id
    };
    LoginActions.submitLoginForm(data, this.context.router);
  },

  /*
   * Clear the state and input and do the focus on the name input
   */
  clearAndFocus() {
    this.focus();
  },

  focus() {
    var n = $(this.refs.name.getDOMNode());
    if (n) {
      n.focus();
    }
  },

  /*
   * If you are already connected, will skip the page
   * listen to the onChange event from the UserStore
   */
  componentDidMount() {
    LoginActions.tryPassLoginForm(this.context.router);

    this.focus();

    if (reactor.evaluateToJS(getters.getIsLogged)) {
      this.replaceWith('/machine');
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

    if (locations && this.state.location) {
    return (
      <form className="login-form" method="post" onSubmit={this.handleSubmit}>
        <div className="regular-login">

          <h2 className="login-heading">Please log in</h2>
          <input
            ref="name"
            type="text"
            name="username"
            className="form-control"
            placeholder="Username"
            required
            autofocus
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
            value={this.state.location.Id}
            onChange={this.updateLocation}
            required>
            {_.map(locations, (location, i) => {
              if (location.Approved) {
                return (
                  <option key={i} value={location.Id}>
                    {location.Title}
                  </option>
                );
              }
            })}
          </select>
          <button className="btn btn-primary btn-block btn-login"
            type="submit">Log In</button>

          {this.state.location ? (
            <div className="signup-link">
              Do not have an account yet? <a href={'/signup/#/form?location=' + this.state.location.Id}
                onClick={this.goToSignUp}>
                Sign up for {this.state.location.Title}
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
    var locationId = parseInt(this.refs.location.getDOMNode().value);
    console.log('location <-', locationId);
    LocationActions.setLocationId(locationId);
  }
});

export default Login;
