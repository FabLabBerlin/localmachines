import React from 'react';

var BusyMachine = React.createClass({
  render() {
    return(
      <div>
        <div className="container-fluid" >
          {this.props.info.name}
          <br/>
          {this.props.activation}
        </div>
        <button className="btn btn-danger" >stop</button>
      </div>
    );
  }
});

module.exports = BusyMachine;
