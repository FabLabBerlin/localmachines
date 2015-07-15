import React from 'react';

var OccupiedMachine = React.createClass({
  render() {
    return (
      <div>
        <div className="container-fluid">
          {this.props.info.Name}
          <br/>
          {this.props.activation}
        </div>
        <div className="indicator indicator-occupied" >occupied</div>
      </div>
    );
  }
});

module.exports = OccupiedMachine;
