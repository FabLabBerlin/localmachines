import React from 'react';

String.prototype.toHHMMSS = function() {
  var d = parseInt(this,10);
  var h = Math.floor(d / 3600);
  var m = Math.floor(d % 3600 / 60);
  var s = Math.floor(d % 3600 % 60);
  console.log(d);
  return ((h > 0 ? h + ':' + (m < 10 ? "0" : "") : "") + m + ":" + (s < 10 ? "0" : "") + s);
}

var BillTable = React.createClass({
  render() {
    if(this.props.info.Details.length != 0) {
      var BillNode = this.props.info.Details.map(function(info) {
        return (
          <tr key={info.MachineId} >
            <td>{info.MachineName}</td>
            <td>{info.Time.toString().toHHMMSS()}</td>
            <td>{(info.Price * 0.19).toFixed(2)} <i className="fa fa-eur"></i></td>
            <td>{info.Price.toFixed(2)} <i className="fa fa-eur"></i></td>
          </tr>
        );
      });
    } else {
      var BillNode = <p>You do not have any expenses</p>
    }
    return (
      <table className="table table-striped table-hover" >
        <thead>
          <tr>
            <th>Machine Name</th>
            <th>Time (h:m:s)</th>
            <th>VAT(19%)</th>
            <th>Expenses <i className="fa fa-eur"></i></th>
          </tr>
        </thead>
        <tbody>
          {BillNode}
          <tr>
            <td><label>Total</label></td>
            <td><label>{this.props.info.TotalTime.toString().toHHMMSS()}</label></td>
            <td><label>{(this.props.info.TotalPrice * 0.19).toFixed(2)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{this.props.info.TotalPrice.toFixed(2)}</label> <i className="fa fa-eur"></i></td>
          </tr>
        </tbody>
      </table>
    );
  }
});

module.exports = BillTable;
