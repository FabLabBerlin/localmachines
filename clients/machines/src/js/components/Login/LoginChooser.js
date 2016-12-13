var Location = require('../../modules/Location');
var Login = require('./Login');
var React = require('react');


/*
 * LoginChooser
 * Will choose the right login page
 * depending of the presence of window.libnfc
 */
var LoginChooser = React.createClass({
  componentWillMount() {
    Location.actions.loadLocations();
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
