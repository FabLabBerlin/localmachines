var _ = require('lodash');
var BillTable = require('./BillTable');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var React = require('react');
var reactor = require('../../reactor');


var BillTables = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      bill: getters.getBill,
      monthlyBills: getters.getMonthlyBills
    };
  },

  render() {
    if (this.state.bill) {
      return (
        <div>
          {_.map(this.state.monthlyBills, (bill, i) => {
            return <BillTable bill={bill} key={i}/>;
          })}
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }
});

export default BillTables;
