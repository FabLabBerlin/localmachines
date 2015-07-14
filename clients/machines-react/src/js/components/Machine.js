import React from 'react';

/*
 * Multiple button available there !!
 * Has to be connected to activation store
 */
var Machine = React.createClass({
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

module.exports = Machine;
