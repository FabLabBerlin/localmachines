var _ = require('lodash');
var getters = require('../../../getters');
var ImageUploader = require('../ImageUploader');
var LoaderLocal = require('../../LoaderLocal');
var LocationActions = require('../../../actions/LocationActions');
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
    LocationActions.loadUserLocations(uid);
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
    const locationId = this.state.locationId;
    console.log('this.state.location=', this.state.location);
    if (!this.state.settings) {
      return <LoaderLocal/>;
    }

    return (
      <div id="set" className="container">
        <h1>Settings</h1>

        <table className="table table-striped">
          <thead/>
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
            <tr>
              <td>Logo</td>
              <td>
                <ImageUploader existingImage={(this.state.location && this.state.location.get('Logo'))
                                               ? '/files/' + this.state.location.get('Logo')
                                               : '/machines/assets/img/logo-easylab.svg'}
                               height="48"
                               uploadUrl={'/api/locations/' + locationId + '/image?location=' + locationId}/>
              </td>
            </tr>
            <tr>
              <td>Reservation Notification E-Mail</td>
              <td>
                <input type="text"
                       ref="ReservationNotificationEmail"
                       defaultValue={this.state.settings.getIn(['ReservationNotificationEmail', 'ValueString'])}/>
              </td>
            </tr>
            <tr>
              <td>Mailchimp API Key</td>
              <td>
                <input type="text"
                       ref="MailchimpApiKey"
                       defaultValue={this.state.settings.getIn(['MailchimpApiKey', 'ValueString'])}/>
              </td>
            </tr>
            <tr>
              <td>Mailchimp List Id</td>
              <td>
                <input type="text"
                       ref="MailchimpListId"
                       defaultValue={this.state.settings.getIn(['MailchimpListId', 'ValueString'])}/>
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
      },
      ReservationNotificationEmail: {
        ValueString: this.refs.ReservationNotificationEmail.value
      },
      MailchimpApiKey: {
        ValueString: this.refs.MailchimpApiKey.value
      },
      MailchimpListId: {
        ValueString: this.refs.MailchimpListId.value
      }
    };

    SettingsActions.saveAdminSettings(settings);
  }

});

export default Settings;
