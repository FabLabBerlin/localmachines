import React from 'react';

var FreeMachine = React.createClass ({
  render() {
    return (
      <div>
        <div className="container-fluid" >
          {this.props.info.Name}
          <br/>
          {this.props.activation}
        </div>
        <button className="btn btn-primary" >start </button>
      </div>
    );
  }
});

module.exports = FreeMachine;
