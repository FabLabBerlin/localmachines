var $ = require('jquery');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var React = require('react');
var reactor = require('../../../reactor');
var SettingsActions = require('../../../modules/Settings/actions');
var SettingsGetters = require('../../../modules/Settings/getters');
var UserActions = require('../../../actions/UserActions');


var Settings = React.createClass({

  mixins: [ reactor.ReactMixin ],

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    SettingsActions.loadAdminSettings({locationId});
    UserActions.fetchUser(uid);
  },

  getDataBindings() {
    return {
      adminSettings: SettingsGetters.getAdminSettings,
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      uid: getters.getUid,
      vatPercent: SettingsGetters.getVatPercent
    };
  },

  render() {
    if (!this.state.adminSettings) {
      return <LoaderLocal/>;
    }

    return (
      <div className="container-fluid">
        <h1>Settings</h1>

        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Setting</th>
              <th>Value</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>URL to AGBs</td>
              <td>
                <input type="text"
                       ref="TermsUrl"
                       defaultValue={this.state.adminSettings.TermsUrl.ValueString}/>
              </td>
            </tr>
            <tr>
              <td>Currency</td>
              <td>
                <input type="text"
                       ref="Currency"
                       defaultValue={this.state.adminSettings.Currency.ValueString}/>
              </td>
            </tr>
            <tr>
              <td>VAT Rate</td>
              <td>
                <input type="number"
                       ref="VAT"
                       defaultValue={this.state.adminSettings.VAT.ValueFloat}/>
              </td>
            </tr>
          </tbody>
        </table>

        <hr/>

        <div className="pull-right">
          <button className="btn btn-primary btn-lg"
                  onClick={this.save}>
            Save
          </button>
        </div>
      </div>
    );
  },

  save() {
    const adminSettings = {
      TermsUrl: {
        ValueString: $(this.refs.TermsUrl.getDOMNode()).val()
      },
      Currency: {
        ValueString: $(this.refs.Currency.getDOMNode()).val()
      },
      VAT: {
        ValueFloat: parseFloat($(this.refs.VAT.getDOMNode()).val())
      }
    };

    console.log('adminSettings=', adminSettings);

    SettingsActions.saveAdminSettings(adminSettings);
  }

});

export default Settings;