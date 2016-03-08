var { Navigation } = require('react-router');
var React = require('react');


var Done = React.createClass({
  mixins: [ Navigation ],

  handleSubmit(event) {
    event.preventDefault();
    this.transitionTo('/login');
  },

  render() {
    return (
      <div className="container text-center">
        <form onSubmit={this.handleSubmit}>
          <h3>You are done!</h3>
          <p>
            Now you can login with your Email address as user name
            and your new password!
          </p>
          <div className="pull-right">
            <button className="btn btn-primary btn-lg wizard-button"
              type="submit">Continue</button>
          </div>
        </form>
      </div>
    );
  }
});

export default Done;