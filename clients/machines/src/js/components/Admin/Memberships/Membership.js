var getters = require('../../../getters');
var {hashHistory} = require('react-router');
var LoaderLocal = require('../../LoaderLocal');
var Location = require('../../../modules/Location');
var Memberships = require('../../../modules/Memberships');
var React = require('react');
var reactor = require('../../../reactor');
var Settings = require('../../../modules/Settings');
var UserActions = require('../../../actions/UserActions');


var MembershipPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      currency: Settings.getters.getCurrency,
      memberships: Memberships.getters.getAllMemberships,
      showArchived: Memberships.getters.getShowArchived
    };
  },

  componentDidMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    Memberships.actions.fetch({locationId});
    Location.actions.loadUserLocations(uid);
    Settings.actions.loadSettings({locationId});
  },

  render() {
    const id = parseInt(this.props.params.membershipId);
    var mb;
    var machine = {};

    if (this.state.memberships) {
      mb = this.state.memberships.find(m => m.get('Id') === id);
    } else {
      return <LoaderLocal/>;
    }

    return (
      <div className="container-fluid">
        <h2>Edit Membership "{mb.get('Title')}"</h2>
        <h3>Basic settings</h3>

        <hr/>

        <div class="row">

          <div class="col-sm-3">
            <div class="form-group">
              <label for="membership-title">Membership Title</label>
              <input type="text" id="membership-title" class="form-control"
                     placeholder="Membership title" value={mb.get('Title')}/>
            </div>
          </div>

          <div class="col-sm-3">
            <div class="form-group">
              <label for="membership-shortname">Shortname</label>
              <input type="text" id="membership-shortname" class="form-control"
                     placeholder="Shortname" value={mb.get('ShortName')}/>
            </div>
          </div>
          
          <div class="col-sm-3">
            <div class="form-group">
              <label for="membership-duration">Duration (months)</label>
              <input type="text" id="membership-duration" class="form-control"
                     placeholder="Duration" value={mb.get('DurationMonths')}/>
            </div>
          </div>
          
          <div class="col-sm-3">
            <div class="form-group">
              <label for="membership-price">Monthly Price</label>
              <div class="input-group">
                <input type="text" id="membership-price" class="form-control"
                       placeholder="Monthly Price" value={mb.get('MonthlyPrice')}/>
                <div class="input-group-addon">
                  {this.state.currency || '€'}
                </div>
              </div>
            </div>
          </div>

        </div>

        <div class="row">
          <div class="col-sm-6">
            <div class="form-group">
              <label>Machines Affected by Membership</label>
              <div class="row">

                <div class="col-sm-6" ng-repeat="machine in machines | machinesFilter:this">
                  <div class="checkbox-inline">
                    <label title={machine.Description}>
                      <input type="checkbox"
                             ng-model="machine['Checked']"/>
                      {machine.Name}
                    </label>
                  </div>
                </div>
                
              </div>
            </div>
          </div>
          <div class="col-sm-3">
            <div class="form-group">
              <label for="membership-price-deduction">Machine price deduction</label>
              <div class="input-group">
                <input type="text" 
                  id="membership-price-deduction" 
                  class="form-control"
                  placeholder="Percentage" 
                  value={mb.get('MachinePriceDeduction')}/>
                <div class="input-group-addon">
                  <b>%</b>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="row">

          <div class="col-sm-3">
            <div class="form-group">
              <label>Automatically extend membership</label>
              <div class="checkbox-inline">
                <label title="Automatically Extend Membership. Or not.">
                  <input type="checkbox" 
                    id="membership-auto-extend"
                    value={mb.get('AutoExtend')}/> 
                    <span ng-show="membership.AutoExtend">Yes</span>
                    <span ng-hide="membership.AutoExtend">No</span>
                </label>
              </div>
            </div>
          </div>

          <div class="col-sm-3">
            <div class="form-group">
              <label for="membership-extend-duration">After the end date membership extends</label>
              <div class="input-group">
                <input type="text"
                  id="membership-extend-duration"
                  class="form-control"
                  value={mb.get('AutoExtendDurationMonths')}
                  ng-disabled="!membership.AutoExtend"/>
                <div class="input-group-addon">-monthly</div>
              </div>
            </div>
          </div>
        
        </div>

        <hr/

        >
      </div>
    );
  }
});

export default MembershipPage;
