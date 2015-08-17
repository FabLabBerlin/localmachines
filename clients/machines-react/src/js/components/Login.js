var React = require('react');
var Flux = require('../flux');
var getters = require('../getters');
var {Navigation} = require('react-router');
var LoginStore = require('../stores/LoginStore');
var LoginActions = require('../actions/LoginActions');


/*
 * Login component
 * Handle the login page
 * Give permission to edit your information
 */
var Login = React.createClass({

  /*
   * To use transitionTo/replaceWith/redirect and some function related to the router
   */
  mixins: [ Navigation ],

  /*
   * Set the initial state
   */
  getInitialState() {
    return {
      username: '',
      password: ''
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
    LoginActions.submitLoginForm(this.state);
    this.clearAndFocus();
  },

  /*
   * Update the state when there are changes in the input
   */
  handleChange(e) {
    this.setState({
      [e.target.name]: e.target.value
    });
  },

  /*
   * Clear the state and input and do the focus on the name input
   */
  clearAndFocus() {
    this.setState({username: '', password: ''}, function() {
      React.findDOMNode(this.refs.name).focus();
    });
  },

  /*
   * Replace the login page url by the user page url
   */
  onChangeLogin() {
    const isLogged = Flux.evaluateToJS(getters.getIsLogged);
    if (isLogged) {
      this.replaceWith('/machine');
    }
  },

  /*
   * If you are already connected, will skip the page
   * listen to the onChange event from the UserStore
   */
  componentDidMount() {
    LoginActions.submitLoginForm(this.state);
    LoginStore.onChangeLogin = this.onChangeLogin;

    Flux.observe(getters.getIsLogged, isLogged => {
      this.onChangeLogin();
    }.bind(this));
  },

  /*
   * Render the form and the button inside of the App component
   */
  render() {
    return (
      <div className="app" >
        <div className="container-fluid">
          <div className="regular-login" >
            <form className="login-form"
              onSubmit={this.handleSubmit} >
              <h2 className="login-heading">Please log in</h2>
              <input
                ref="name"
                type="text"
                name="username"
                className="form-control"
                value={this.state.username}
                onChange={this.handleChange}
                placeholder="Username"
                required
                autofocus
              />
              <input
                type="password"
                name="password"
                className="form-control"
                ref="Password"
                value={this.state.password}
                onChange={this.handleChange}
                placeholder="password"
                required
              />
              <button className="btn btn-primary btn-block btn-login"
                type="submit">Log In</button>

              <div className="signup-link">
                <button className="btn btn-lg btn-info"
                  onClick={this.goToSignUp}>
                  <i className="fa fa-sign-in"></i>
                  Sign up
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
});

export default Login;
