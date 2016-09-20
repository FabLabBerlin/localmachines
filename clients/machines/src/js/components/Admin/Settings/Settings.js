var _ = require('lodash');
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
    SettingsActions.loadSettings({locationId});
    SettingsActions.loadFastbillTemplates({locationId});
    UserActions.fetchUser(uid);
  },

  getDataBindings() {
    return {
      settings: SettingsGetters.getAdminSettings,
      fastbillTemplates: SettingsGetters.getFastbillTemplates,
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      uid: getters.getUid,
      vatPercent: SettingsGetters.getVatPercent
    };
  },

  render() {
    if (!this.state.settings) {
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
                       defaultValue={this.state.settings.getIn(['TermsUrl', 'ValueString'])}/>
              </td>
            </tr>
            <tr>
              <td>Currency</td>
              <td>
                <input type="text"
                       ref="Currency"
                       defaultValue={this.state.settings.getIn(['Currency', 'ValueString'])}/>
              </td>
            </tr>
            <tr>
              <td>VAT Rate</td>
              <td>
                <input type="number"
                       ref="VAT"
                       defaultValue={this.state.settings.getIn(['VAT', 'ValueFloat'])}/>
              </td>
            </tr>
            <tr>
              <td>Fastbill Template</td>
              <td>
                {this.state.fastbillTemplates ?
                  (
                    <select ref="FastbillTemplateId"
                            defaultValue={this.state.settings.getIn(['FastbillTemplateId', 'ValueInt'])}>
                      <option value="0">Please select</option>
                      {_.map(this.state.fastbillTemplates.toJS(), (t) => {
                        return (
                          <option key={t.TEMPLATE_ID}
                                  value={t.TEMPLATE_ID}>
                            {t.TEMPLATE_NAME}
                          </option>
                        );
                      })}
                    </select>
                  ) : <LoaderLocal/>}
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
    const settings = {
      TermsUrl: {
        ValueString: this.refs.TermsUrl.value
      },
      Currency: {
        ValueString: this.refs.Currency.value
      },
      FastbillTemplateId: {
        ValueInt: parseInt(this.refs.FastbillTemplateId.value)
      },
      VAT: {
        ValueFloat: parseFloat(this.refs.VAT.value)
      }
    };

    SettingsActions.saveAdminSettings(settings);
  }

});

export default Settings;
