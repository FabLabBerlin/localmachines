import React from 'react';

var BillTable = React.createClass({
  render() {
    if(this.props.info.Details.length != 0) {
      var BillNode = this.props.info.Details.map(function(info) {
        return (
          <tr key={info.MachineId} >
            <td>{info.MachineName}</td>
            <td>{info.Time}</td>
            <td>{info.Price} <i className="fa fa-eur"></i></td>
          </tr>
        );
      });
    } else {
      var BillNode = <p>You do not have any credit</p>
    }
    return (
      <table className="table table-striped table-hover" >
        <thead>
          <tr>
            <th>Machine Name</th>
            <th>Time (s)</th>
            <th>Credit <i className="fa fa-eur"></i></th>
          </tr>
        </thead>
        <tbody>
          {BillNode}
          <tr>
            <td><label>Total</label></td>
            <td><label>{this.props.info.TotalTime}</label></td>
            <td><label>{this.props.info.TotalPrice}</label> <i className="fa fa-eur"></i></td>
          </tr>
        </tbody>
      </table>
    );
  }
});

module.exports = BillTable;
