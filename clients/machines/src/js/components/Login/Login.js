var $ = require('jquery');

var getters = require('../../getters');

var LoginStore = require('../../stores/LoginStore');
var LocationStore = require('../../stores/LocationStore');

var LoginActions = require('../../actions/LoginActions');

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
      locations: getters.getLocations
    };
  },

  goToSignUp(event) {
    event.preventDefault();
    window.location = '/signup';
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
      location: this.refs.location.getDOMNode().value
    };
    console.log('handleSubmit: username=', data.username);
    console.log('handleSubmit: location=', data.location);
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
            placeholder="password"
            required
          />
          <select
            className="form-control location-picker"
            ref="location"
            name="location"
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

          <div className="signup-link">
            Do not have an account yet? <a href="#"
              onClick={this.goToSignUp}>
              Sign up
            </a> now!
          </div>

        </div>
      </form>
    );
    } else {
      return (<div></div>);
    }
  }
});

export default Login;
