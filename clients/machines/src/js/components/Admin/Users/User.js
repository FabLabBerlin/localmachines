var _ = require('lodash');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationActions = require('../../../actions/LocationActions');
var LocationGetters = require('../../../modules/Location/getters');
var MachineActions = require('../../../actions/MachineActions');
var Machines = require('../../../modules/Machines');
var Navigation = require('react-router').Navigation;
var React = require('react');
var reactor = require('../../../reactor');
var UserActions = require('../../../actions/UserActions');
var Users = require('../../../modules/Users');


const countryCodesJSON = '[{"Name":"Afghanistan","Code":"AF"},{"Name":"Åland Islands","Code":"AX"},{"Name":"Albania","Code":"AL"},{"Name":"Algeria","Code":"DZ"},{"Name":"American Samoa","Code":"AS"},{"Name":"Andorra","Code":"AD"},{"Name":"Angola","Code":"AO"},{"Name":"Anguilla","Code":"AI"},{"Name":"Antarctica","Code":"AQ"},{"Name":"Antigua and Barbuda","Code":"AG"},{"Name":"Argentina","Code":"AR"},{"Name":"Armenia","Code":"AM"},{"Name":"Aruba","Code":"AW"},{"Name":"Australia","Code":"AU"},{"Name":"Austria","Code":"AT"},{"Name":"Azerbaijan","Code":"AZ"},{"Name":"Bahamas","Code":"BS"},{"Name":"Bahrain","Code":"BH"},{"Name":"Bangladesh","Code":"BD"},{"Name":"Barbados","Code":"BB"},{"Name":"Belarus","Code":"BY"},{"Name":"Belgium","Code":"BE"},{"Name":"Belize","Code":"BZ"},{"Name":"Benin","Code":"BJ"},{"Name":"Bermuda","Code":"BM"},{"Name":"Bhutan","Code":"BT"},{"Name":"Bolivia, Plurinational State of","Code":"BO"},{"Name":"Bonaire, Sint Eustatius and Saba","Code":"BQ"},{"Name":"Bosnia and Herzegovina","Code":"BA"},{"Name":"Botswana","Code":"BW"},{"Name":"Bouvet Island","Code":"BV"},{"Name":"Brazil","Code":"BR"},{"Name":"British Indian Ocean Territory","Code":"IO"},{"Name":"Brunei Darussalam","Code":"BN"},{"Name":"Bulgaria","Code":"BG"},{"Name":"Burkina Faso","Code":"BF"},{"Name":"Burundi","Code":"BI"},{"Name":"Cambodia","Code":"KH"},{"Name":"Cameroon","Code":"CM"},{"Name":"Canada","Code":"CA"},{"Name":"Cape Verde","Code":"CV"},{"Name":"Cayman Islands","Code":"KY"},{"Name":"Central African Republic","Code":"CF"},{"Name":"Chad","Code":"TD"},{"Name":"Chile","Code":"CL"},{"Name":"China","Code":"CN"},{"Name":"Christmas Island","Code":"CX"},{"Name":"Cocos (Keeling) Islands","Code":"CC"},{"Name":"Colombia","Code":"CO"},{"Name":"Comoros","Code":"KM"},{"Name":"Congo","Code":"CG"},{"Name":"Congo, the Democratic Republic of the","Code":"CD"},{"Name":"Cook Islands","Code":"CK"},{"Name":"Costa Rica","Code":"CR"},{"Name":"Côte d\'Ivoire","Code":"CI"},{"Name":"Croatia","Code":"HR"},{"Name":"Cuba","Code":"CU"},{"Name":"Curaçao","Code":"CW"},{"Name":"Cyprus","Code":"CY"},{"Name":"Czech Republic","Code":"CZ"},{"Name":"Denmark","Code":"DK"},{"Name":"Djibouti","Code":"DJ"},{"Name":"Dominica","Code":"DM"},{"Name":"Dominican Republic","Code":"DO"},{"Name":"Ecuador","Code":"EC"},{"Name":"Egypt","Code":"EG"},{"Name":"El Salvador","Code":"SV"},{"Name":"Equatorial Guinea","Code":"GQ"},{"Name":"Eritrea","Code":"ER"},{"Name":"Estonia","Code":"EE"},{"Name":"Ethiopia","Code":"ET"},{"Name":"Falkland Islands (Malvinas)","Code":"FK"},{"Name":"Faroe Islands","Code":"FO"},{"Name":"Fiji","Code":"FJ"},{"Name":"Finland","Code":"FI"},{"Name":"France","Code":"FR"},{"Name":"French Guiana","Code":"GF"},{"Name":"French Polynesia","Code":"PF"},{"Name":"French Southern Territories","Code":"TF"},{"Name":"Gabon","Code":"GA"},{"Name":"Gambia","Code":"GM"},{"Name":"Georgia","Code":"GE"},{"Name":"Germany","Code":"DE"},{"Name":"Ghana","Code":"GH"},{"Name":"Gibraltar","Code":"GI"},{"Name":"Greece","Code":"GR"},{"Name":"Greenland","Code":"GL"},{"Name":"Grenada","Code":"GD"},{"Name":"Guadeloupe","Code":"GP"},{"Name":"Guam","Code":"GU"},{"Name":"Guatemala","Code":"GT"},{"Name":"Guernsey","Code":"GG"},{"Name":"Guinea","Code":"GN"},{"Name":"Guinea-Bissau","Code":"GW"},{"Name":"Guyana","Code":"GY"},{"Name":"Haiti","Code":"HT"},{"Name":"Heard Island and McDonald Islands","Code":"HM"},{"Name":"Holy See (Vatican City State)","Code":"VA"},{"Name":"Honduras","Code":"HN"},{"Name":"Hong Kong","Code":"HK"},{"Name":"Hungary","Code":"HU"},{"Name":"Iceland","Code":"IS"},{"Name":"India","Code":"IN"},{"Name":"Indonesia","Code":"ID"},{"Name":"Iran, Islamic Republic of","Code":"IR"},{"Name":"Iraq","Code":"IQ"},{"Name":"Ireland","Code":"IE"},{"Name":"Isle of Man","Code":"IM"},{"Name":"Israel","Code":"IL"},{"Name":"Italy","Code":"IT"},{"Name":"Jamaica","Code":"JM"},{"Name":"Japan","Code":"JP"},{"Name":"Jersey","Code":"JE"},{"Name":"Jordan","Code":"JO"},{"Name":"Kazakhstan","Code":"KZ"},{"Name":"Kenya","Code":"KE"},{"Name":"Kiribati","Code":"KI"},{"Name":"Korea, Democratic People\'s Republic of","Code":"KP"},{"Name":"Korea, Republic of","Code":"KR"},{"Name":"Kuwait","Code":"KW"},{"Name":"Kyrgyzstan","Code":"KG"},{"Name":"Lao People\'s Democratic Republic","Code":"LA"},{"Name":"Latvia","Code":"LV"},{"Name":"Lebanon","Code":"LB"},{"Name":"Lesotho","Code":"LS"},{"Name":"Liberia","Code":"LR"},{"Name":"Libya","Code":"LY"},{"Name":"Liechtenstein","Code":"LI"},{"Name":"Lithuania","Code":"LT"},{"Name":"Luxembourg","Code":"LU"},{"Name":"Macao","Code":"MO"},{"Name":"Macedonia, the Former Yugoslav Republic of","Code":"MK"},{"Name":"Madagascar","Code":"MG"},{"Name":"Malawi","Code":"MW"},{"Name":"Malaysia","Code":"MY"},{"Name":"Maldives","Code":"MV"},{"Name":"Mali","Code":"ML"},{"Name":"Malta","Code":"MT"},{"Name":"Marshall Islands","Code":"MH"},{"Name":"Martinique","Code":"MQ"},{"Name":"Mauritania","Code":"MR"},{"Name":"Mauritius","Code":"MU"},{"Name":"Mayotte","Code":"YT"},{"Name":"Mexico","Code":"MX"},{"Name":"Micronesia, Federated States of","Code":"FM"},{"Name":"Moldova, Republic of","Code":"MD"},{"Name":"Monaco","Code":"MC"},{"Name":"Mongolia","Code":"MN"},{"Name":"Montenegro","Code":"ME"},{"Name":"Montserrat","Code":"MS"},{"Name":"Morocco","Code":"MA"},{"Name":"Mozambique","Code":"MZ"},{"Name":"Myanmar","Code":"MM"},{"Name":"Namibia","Code":"NA"},{"Name":"Nauru","Code":"NR"},{"Name":"Nepal","Code":"NP"},{"Name":"Netherlands","Code":"NL"},{"Name":"New Caledonia","Code":"NC"},{"Name":"New Zealand","Code":"NZ"},{"Name":"Nicaragua","Code":"NI"},{"Name":"Niger","Code":"NE"},{"Name":"Nigeria","Code":"NG"},{"Name":"Niue","Code":"NU"},{"Name":"Norfolk Island","Code":"NF"},{"Name":"Northern Mariana Islands","Code":"MP"},{"Name":"Norway","Code":"NO"},{"Name":"Oman","Code":"OM"},{"Name":"Pakistan","Code":"PK"},{"Name":"Palau","Code":"PW"},{"Name":"Palestine, State of","Code":"PS"},{"Name":"Panama","Code":"PA"},{"Name":"Papua New Guinea","Code":"PG"},{"Name":"Paraguay","Code":"PY"},{"Name":"Peru","Code":"PE"},{"Name":"Philippines","Code":"PH"},{"Name":"Pitcairn","Code":"PN"},{"Name":"Poland","Code":"PL"},{"Name":"Portugal","Code":"PT"},{"Name":"Puerto Rico","Code":"PR"},{"Name":"Qatar","Code":"QA"},{"Name":"Réunion","Code":"RE"},{"Name":"Romania","Code":"RO"},{"Name":"Russian Federation","Code":"RU"},{"Name":"Rwanda","Code":"RW"},{"Name":"Saint Barthélemy","Code":"BL"},{"Name":"Saint Helena, Ascension and Tristan da Cunha","Code":"SH"},{"Name":"Saint Kitts and Nevis","Code":"KN"},{"Name":"Saint Lucia","Code":"LC"},{"Name":"Saint Martin (French part)","Code":"MF"},{"Name":"Saint Pierre and Miquelon","Code":"PM"},{"Name":"Saint Vincent and the Grenadines","Code":"VC"},{"Name":"Samoa","Code":"WS"},{"Name":"San Marino","Code":"SM"},{"Name":"Sao Tome and Principe","Code":"ST"},{"Name":"Saudi Arabia","Code":"SA"},{"Name":"Senegal","Code":"SN"},{"Name":"Serbia","Code":"RS"},{"Name":"Seychelles","Code":"SC"},{"Name":"Sierra Leone","Code":"SL"},{"Name":"Singapore","Code":"SG"},{"Name":"Sint Maarten (Dutch part)","Code":"SX"},{"Name":"Slovakia","Code":"SK"},{"Name":"Slovenia","Code":"SI"},{"Name":"Solomon Islands","Code":"SB"},{"Name":"Somalia","Code":"SO"},{"Name":"South Africa","Code":"ZA"},{"Name":"South Georgia and the South Sandwich Islands","Code":"GS"},{"Name":"South Sudan","Code":"SS"},{"Name":"Spain","Code":"ES"},{"Name":"Sri Lanka","Code":"LK"},{"Name":"Sudan","Code":"SD"},{"Name":"Suriname","Code":"SR"},{"Name":"Svalbard and Jan Mayen","Code":"SJ"},{"Name":"Swaziland","Code":"SZ"},{"Name":"Sweden","Code":"SE"},{"Name":"Switzerland","Code":"CH"},{"Name":"Syrian Arab Republic","Code":"SY"},{"Name":"Taiwan, Province of China","Code":"TW"},{"Name":"Tajikistan","Code":"TJ"},{"Name":"Tanzania, United Republic of","Code":"TZ"},{"Name":"Thailand","Code":"TH"},{"Name":"Timor-Leste","Code":"TL"},{"Name":"Togo","Code":"TG"},{"Name":"Tokelau","Code":"TK"},{"Name":"Tonga","Code":"TO"},{"Name":"Trinidad and Tobago","Code":"TT"},{"Name":"Tunisia","Code":"TN"},{"Name":"Turkey","Code":"TR"},{"Name":"Turkmenistan","Code":"TM"},{"Name":"Turks and Caicos Islands","Code":"TC"},{"Name":"Tuvalu","Code":"TV"},{"Name":"Uganda","Code":"UG"},{"Name":"Ukraine","Code":"UA"},{"Name":"United Arab Emirates","Code":"AE"},{"Name":"United Kingdom","Code":"GB"},{"Name":"United States","Code":"US"},{"Name":"United States Minor Outlying Islands","Code":"UM"},{"Name":"Uruguay","Code":"UY"},{"Name":"Uzbekistan","Code":"UZ"},{"Name":"Vanuatu","Code":"VU"},{"Name":"Venezuela, Bolivarian Republic of","Code":"VE"},{"Name":"Viet Nam","Code":"VN"},{"Name":"Virgin Islands, British","Code":"VG"},{"Name":"Virgin Islands, U.S.","Code":"VI"},{"Name":"Wallis and Futuna","Code":"WF"},{"Name":"Western Sahara","Code":"EH"},{"Name":"Yemen","Code":"YE"},{"Name":"Zambia","Code":"ZM"},{"Name":"Zimbabwe","Code":"ZW"}]';
const countryCodes = JSON.parse(countryCodesJSON);


var GeneralInfo = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    };
  },

  render() {
    const user = this.props.user;

    return (
      <div>
        <div className="row">

          <div className="col-sm-3">
            <div className="form-group">
              <label for="user-first-name">First Name</label>
              <input id="user-first-name" type="text" className="form-control" 
                     placeholder="Enter first name" defaultValue={user.FirstName}/>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label for="user-last-name">Last Name</label>
              <input id="user-last-name" type="text" className="form-control" 
                     placeholder="Enter last name" defaultValue={user.LastName}/>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label for="user-username">Username</label>
              <input id="user-username" type="text" className="form-control" 
                     placeholder="Enter Username" defaultValue={user.Username}/>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label for="user-email">E-Mail</label>
              <input id="user-email" type="text" className="form-control"
                     placeholder="Enter E-Mail" defaultValue={user.Email}/>
            </div>
          </div>

        </div>

        <div className="row">
          <div className="col-sm-3">
            <label for="phone">Phone</label>
            <input 
              type="text" 
              className="form-control"
              id="phone"
              defaultValue={user.Phone}/>
          </div>
          <div className="col-sm-3">
            <div className="form-group">
              <label for="invoice-addr">Billing Address</label>
              <input id="invoice-addr"
                className="form-control"
                placeholder="Billing Address"
                defaultValue={user.InvoiceAddr}/>
            </div>
          </div>
          <div className="col-sm-3">
            <div className="form-group">
              <label for="zip-code">ZIP Code</label>
              <input id="zip-code"
                className="form-control"
                placeholder="ZIP"
                defaultValue={user.ZipCode}/>
            </div>
          </div>
          <div className="col-sm-3">
            <div className="form-group">
              <label for="city">City</label>
              <input id="city"
                className="form-control"
                placeholder="City"
                defaultValue={user.City}/>
            </div>
          </div>
        </div>

        <div className="row">
          <div className="col-sm-3">
            <div className="form-group">
              <label for="country">Country</label>
              <select id="country"
                      className="form-control" 
                      value={user.CountryCode}>
                <option value="" selected disabled>Select Country</option>
                {_.each(countryCodes, (country) => {
                  return (
                    <option value={country.Code}>
                      {country.Name}
                    </option>
                  );
                })}
              </select>
            </div>
          </div>
          <div className="col-sm-3">
            <div className="form-group">
              <label for="company">Company</label>
              <input id="company"
                className="form-control"
                placeholder="Company"
                defaultValue={user.Company}/>
            </div>
          </div>
        </div>
      </div>
    );
  }

});


var BillingInfo = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    };
  },

  render() {
    const user = this.props.user;

    return (
      <div>
        <div className="row">
          <div className="col-sm-3">
            <div className="form-group">
              <label for="fastbill-customer-number">FastBill Customer Number</label>
              <input 
                id="fastbill-customer-number"
                className="form-control"
                value={user.ClientId}/>
            </div>
          </div>
          <div className="col-sm-3">
            <div className="form-group">
              <label for="fastbill-action">FastBill Action</label>
              <select id="fastbill-action"
                      className="form-control">
                <option value="1">Create Fastbill customer (gets customer #)</option>
                <option value="2">Load data from FastBill</option>
                <option value="3">Update FastBill data</option>
              </select>
            </div>
          </div>
          <div className="col-sm-3">
            <div className="form-group">
              <label>Sync</label>
              <button 
                ng-click="syncWithFastBill()"
                className="btn btn-primary btn-block">
                Sync with FastBill
              </button>  
            </div>
          </div>
          <div className="col-sm-3">
          </div>
        </div>

        <div className="row">
          <div className="col-sm-3">
            <div className="form-group">
              <label>
                <input type="checkbox" ng-model="user.NoAutoInvoicing"/> No Automatic Invoicing
              </label>
            </div>
          </div>
        </div>
      </div>
    );
  }

});


var Password = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    };
  },

  render() {
    return (
      <div>
        <div>
          <label for="user-password">User Password</label>
        </div>

        <div className="row">
          
            <div className="col-sm-3">
              <div className="form-group">
                <input type="password" className="form-control" 
                       id="user-password" placeholder="New password"/>
              </div>
            </div>

            <div className="col-sm-3">
              <div className="form-group">
                <input type="password" className="form-control"
                       placeholder="Repeat password"/>
              </div>
            </div>

            <div className="col-sm-3">
              <div className="form-group">
                <button className="btn btn-primary btn-block" ng-click="updatePassword()">
                  Update Password
                </button>
              </div>
            </div>
          
        </div>
      </div>
    );
  }

});


var Comments = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    };
  },

  render() {
    const user = this.props.user;

    return (
      <div className="row">
        <div className="col-sm-12">
          <h2>Comments</h2>
        </div>
        <div className="col-sm-6">
          <textarea className="form-control" value={user.Comments}></textarea>
        </div>
      </div>
    );
  }

});


var Permissions = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      machines: Machines.getters.getMachines,
      userLocation: LocationGetters.getUserLocation
    };
  },

  render() {
    if (!this.state.userLocation || !this.state.machines) {
      return <LoaderLocal/>;
    }

    const userLocation = this.state.userLocation.toJS();
    const machines = this.state.machines.toJS();

    console.log('Permissions: machines=', machines);

    return (
      <div>
        <div className="row">

          <div className="col-sm-6">
            <div className="row">
              <div className="col-sm-12">
                <h2>Role</h2>
              </div>
              <div className="col-sm-12 form-group">
                <div className="col-sm-12" ng-repeat="userLocation in userLocations">
                  <div ng-show="globalConfigVisible || userLocation.LocationId == locationId">
                    <div className="col-sm-6">
                      {userLocation.Location.Title}
                    </div>
                    <div className="col-sm-6">
                      <select className="form-control"
                              value={userLocation.UserRole}
                              ng-change="updateUserLocation(userLocation)">
                        <option value="archived">Archived</option>
                        <option value="member">Member</option>
                        <option value="staff">Staff</option>
                        <option value="api">Api</option>
                        <option value="admin">Admin</option>
                      </select>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="col-sm-6">
            <div className="row">
              
              <div className="col-sm-12">
                <h2>Machine Permissions</h2>
              </div>

              <div className="col-sm-6" ng-repeat="machine in availableMachines | machinesFilter:this">
                {_.map(machines, (machine) => {
                  if (machine.Archived) {
                    return undefined;
                  }

                  return (
                    <div className="checkbox-inline">
                      <label title={machine.Description}>
                        <input 
                          type="checkbox" 
                          ng-model="machine['Checked']"
                          ng-disabled="machine['Disabled']"/> 
                        {machine.Name}
                      </label>
                    </div>
                  );
                })}
              </div>
              
            </div>
          </div>

        </div>
      </div>
    );
  }

});


var UserMemberships = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    };
  },

  render() {
    var membership = {};
    var userMembership = {};

    return (
      <div>
        <h2>User Memberships</h2>
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Membership Name</th>
              <th>Start Date</th>
              <th>End Date</th>
              <th>Extends Automatically</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr ng-hide="userMemberships.length">
              <td colspan="6">User has no memberships</td>
            </tr>
            <tr ng-repeat="userMembership in userMemberships">
              <td>{userMembership.Title}</td>
              <td style="min-width:90px">{userMembership.StartDate}</td>
              <td style="min-width:160px">
                <div className="form-inline">
                  <div className="input-group">
                    <input
                      className="adm-user-membership-end-date form-control datepicker"
                      data-user-membership-id="{userMembership.Id}"
                      placeholder="End Date"
                      value="{userMembership.EndDate}"/>
                    <div className="input-group-addon">
                      <i className="fa fa-calendar"/>
                    </div>
                  </div>
                </div>
              </td>
              <td>
                <input 
                  type="checkbox" 
                  ng-model="userMembership.AutoExtend"
                  ng-disabled="userMemberhip.Inactive"
                  ng-change="updateUserMembership(userMembership.Id)"/>
              </td>
              <td>
                <div ng-show="userMembership.Active">Active</div>
                <div ng-show="userMembership.Cancelled">Cancelled</div>
                <div ng-show="userMembership.Inactive">Inactive</div>
              </td>
            </tr>
          </tbody>
        </table>

        <div className="row">
          <div className="col-sm-3">
            <div className="form-group">
              <select 
                className="form-control" 
                id="user-select-membership" 
                placeholder="Membership">
                <option value="" selected disabled>Select Membership</option>
                <option ng-repeat="membership in memberships | membershipsFilter:this" 
                        value={membership.Id}>
                  {membership.Title}
                </option>
              </select>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <div className="input-group">
                <input 
                  id="adm-add-user-membership-start-date" 
                  className="form-control datepicker" 
                  placeholder="Start Date"/>
                <div className="input-group-addon">
                  <i className="fa fa-calendar"></i>
                </div>
              </div>
            </div>
          </div>

          <div className="col-sm-3">
            <button 
              className="btn btn-primary btn-block" 
              id="user-add-membership-btn" 
              ng-click="addUserMembership()">
              <i className="fa fa-plus"></i>&nbsp;Add Membership
            </button>
          </div>
        </div>
      </div>
    );
  }

});


var Buttons = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    };
  },

  render() {
    return (
      <div>
      </div>
    );
  }

});


var UserView = React.createClass({

  mixins: [ Navigation, reactor.ReactMixin ],

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    Users.actions.fetchUsers({locationId});
    UserActions.fetchUser(uid);
    LocationActions.loadUserLocations(uid);
  },

  getDataBindings() {
    return {
      locations: LocationGetters.getLocations,
      location: LocationGetters.getLocation,
      users: Users.getters.getUsers
    };
  },

  render() {
    console.log('UserView: this.state.locations=', this.state.locations);
    if (this.state.users && this.state.locations) {
      const userId = parseInt(this.props.params.userId);
      const user = _.find(this.state.users, u => u.Id === userId);

      return (
        <div className="container-fluid">
          <h1>Edit User</h1>
          <p>Created on {user.Created}</p>

          <hr/>

          <GeneralInfo user={user}/>
          <BillingInfo user={user}/>
          <Password user={user}/>
          <Comments user={user}/>
          <Permissions user={user}/>
          <UserMemberships user={user}/>
          <Buttons user={user}/>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }
});

export default UserView;
