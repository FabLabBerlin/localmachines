import Button from '../../Button';
import getters from '../../../getters';
import Invoices from '../../../modules/Invoices';
import List from './List';
import LoaderLocal from '../../LoaderLocal';
import Location from '../../../modules/Location';
import moment from 'moment';
import React from 'react';
import reactor from '../../../reactor';
import Settings from '../../../modules/Settings';
import SettingsActions from '../../../modules/Settings/actions';
import UserActions from '../../../actions/UserActions';


var ToggleInactiveUsers = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      showInactiveUsers: Invoices.getters.getShowInactiveUsers
    };
  },

  render() {
    if (this.state.showInactiveUsers) {
      return <Button.Annotated id="invs-toggle-inactive"
                               icon="/machines/assets/img/invoicing/inactive_hide.svg"
                               label="Hide inactive users"
                               onClick={this.setShow.bind(this, false)}/>;
    } else {
      return <Button.Annotated id="invs-toggle-inactive"
                               icon="/machines/assets/img/invoicing/inactive_show.svg"
                               label="Show inactive users"
                               onClick={this.setShow.bind(this, true)}/>;
    }
  },

  setShow(yes, e) {
    e.stopPropagation();
    Invoices.actions.setShowInactiveUsers(yes);
  }
});


var Month = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      currency: Settings.getters.getCurrency,
      location: Location.getters.getLocation,
      locationId: Location.getters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      uid: getters.getUid
    };
  },

  componentWillMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);
  },

  checkedComplete(e) {
    e.stopPropagation();
    Invoices.actions.checkedComplete(this.state.locationId);
  },

  checkedPushDrafts(e) {
    e.stopPropagation();
    Invoices.actions.checkedPushDrafts(this.state.locationId);
  },

  checkedSend(e) {
    e.stopPropagation();
    Invoices.actions.checkedSend(this.state.locationId);
  },

  isSelected() {
    const t = this.props.t;
    const month = this.state.MonthlySums
                      .get('selected').get('month');
    const year = this.state.MonthlySums
                      .get('selected').get('year');
    return month === t.month() + 1 &&
           year === t.year();
  },

  render() {
    const t = this.props.t;
    const month = this.state.MonthlySums
                      .get('selected').get('month');
    const year = this.state.MonthlySums
                      .get('selected').get('year');
    const selected = this.isSelected();
    const summaries = this.state.MonthlySums.getIn([year, month]);

    if (!summaries && selected) {
      return <LoaderLocal/>;
    }

    var total;
    if (selected) {
      total = (Math.round(this.state.MonthlySums.getIn([
        t.year(),
        t.month() + 1
      ]).filter(monthlySum => {
        return !monthlySum.get('Canceled');
      }).reduce((result, monthlySum) => {
        return result + monthlySum.get('Total');
      }, 0) * 100) / 100).toFixed(2);
    }

    return (
      <div className={'inv-monthly-sums ' + (selected ? 'selected' : '')}>
        {selected ?
          (
            <div className="row" onClick={this.select}>
              <div className="col-md-6">
                <ToggleInactiveUsers/>
              </div>
              <div className="col-md-6 text-md-right">
                <button type="button"
                        onClick={this.checkedSend}
                        title="Send Invoices">
                  <img src="/machines/assets/img/invoicing/send_invoice_white.svg"/>
                </button>
                <button type="button"
                        onClick={this.checkedComplete}
                        title="Create Invoices">
                  <i className="fa fa-money"/>
                </button>
                <button type="button"
                        onClick={this.checkedPushDrafts}
                        title="Upload Invoice Drafts to Fastbill">
                  <i className="fa fa-cloud-upload"/>
                </button>
              </div>
            </div>
          ) : null
        }
        <div className="row" onClick={this.select}>
          <div className="col-sm-6">
            <h3>{t.format('MMMM YYYY')}</h3>
          </div>
          <div className="col-sm-6 text-sm-right">
            <h3>
              {selected ? ('Sum total: ' + total + ' ' + this.state.currency) : null}
            </h3>
          </div>
        </div>
        {
          selected ? (
            summaries ? (
              <List summaries={summaries}/>
            ) : (
              <LoaderLocal/>
            )
          ) : null
        }
      </div>
    );
  },

  select() {
    if (this.isSelected()) {
      Invoices.actions.setSelectedMonth({
        month: undefined,
        year: undefined
      });
    } else {
      const t = this.props.t;
      const summaries = this.state.MonthlySums.getIn([t.year(), t.month() + 1]);

      Invoices.actions.setSelectedMonth({
        month: t.month() + 1,
        year: t.year()
      });

      if (!summaries) {
        Invoices.actions.fetchMonthlySums(this.state.locationId, {
          month: t.month() + 1,
          year: t.year()
        });
      }
    }
  }

});


var InvoicesView = React.createClass({

  mixins: [ reactor.ReactMixin ],

  componentWillMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    var t = moment();
    Invoices.actions.fetchMonthlySums(this.state.locationId, {
      month: t.month() + 1,
      year: t.year()
    });
    SettingsActions.loadSettings({locationId});
  },

  getDataBindings() {
    return {
      location: Location.getters.getLocation,
      locationId: Location.getters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums
    };
  },

  render() {
    var t = moment();
    var nodes = [];

    for (var i = 0; i < 12; i++) {
      nodes.push(
        <Month key={i} t={t}/>
      );
      t = t.clone().subtract(1, 'months');
    }
    return (
      <div className="container-fluid">
        <h2>Invoices</h2>
        {nodes}
      </div>
    );
  }

});

export default InvoicesView;
