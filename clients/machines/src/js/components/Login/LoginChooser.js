var LocationActions = require('../../actions/LocationActions');
var Login = require('./Login');
var React = require('react');


/*
 * LoginChooser
 * Will choose the right login page
 * depending of the presence of window.libnfc
 */
var LoginChooser = React.createClass({
  componentWillMount() {
    LocationActions.loadLocations();
  },

  render() {
    return (
      <div className="login">
        <div className="container-fluid">

          <Login />

        </div>
      </div>
    );
  }
});

export default LoginChooser;
