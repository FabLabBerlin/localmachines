import _ from 'lodash';
import BillTable from './BillTable';
import getters from '../../getters';
import LoaderLocal from '../LoaderLocal';
import React from 'react';
import reactor from '../../reactor';


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
      var nodes = [];

      this.state.monthlyBills.forEach((bill, i) => {
        nodes.push(<BillTable invoice={bill} key={i}/>);
      });

      return (
        <div>
          {nodes}
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }
});

export default BillTables;
