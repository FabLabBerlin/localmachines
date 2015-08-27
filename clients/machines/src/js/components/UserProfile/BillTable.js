var React = require('react');


function formatTime(t) {
  if (t) {
    var d = parseInt(t.toString(), 10);
    var h = Math.floor(d / 3600);
    var m = Math.floor(d % 3600 / 60);
    var s = Math.floor(d % 3600 % 60);
    var str = '';
    if (h) {
      str += String(h) + ' h ';
    }
    if (h || m) {
      str += String(m) + ' m ';
    }
    if (s) {
      str += String(s) + ' s ';
    }
    return str;
  }
}

var BillTable = React.createClass({
  render() {
    var BillNode;
    var VAT = 0.19;
    if (this.props.info.Details && this.props.info.Details.length !== 0) {
      BillNode = this.props.info.Details.map(function(info) {
        return (
          <tr key={info.MachineId} >
            <td>{info.MachineName}</td>
            <td>{formatTime(info.Time)}</td>
            <td>{(info.Price * VAT).toFixed(2)} <i className="fa fa-eur"></i></td>
            <td>{info.Price.toFixed(2)} <i className="fa fa-eur"></i></td>
          </tr>
        );
      });
    } else {
      return <p>You do not have any expenses</p>;
    }
    return (
      <table className="table table-striped table-hover" >
        <thead>
          <tr>
            <th>Machine</th>
            <th>Time</th>
            <th>VAT (19%)</th>
            <th>Expenses <i className="fa fa-eur"></i> (VAT included)</th>
          </tr>
        </thead>
        <tbody>
          {BillNode}
          <tr>
            <td><label>Total</label></td>
            <td><label>{formatTime(this.props.info.TotalTime)}</label></td>
            <td><label>{(this.props.info.TotalPrice * VAT).toFixed(2)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{this.props.info.TotalPrice.toFixed(2)}</label> <i className="fa fa-eur"></i></td>
          </tr>
        </tbody>
      </table>
    );
  }
});

module.exports = BillTable;
