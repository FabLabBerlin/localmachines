var React = require('react');


var EmailSent = React.createClass({
  render() {
    return (
      <div className="container text-center">
        <h3>Check your mailbox!</h3>
        <p>
          Password reset mail sent to your E-Mail address
        </p>
      </div>
    );
  }
});

export default EmailSent;
