var Categories = require('../../../modules/Categories');
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
      categories: Categories.getters.getAll,
      currency: Settings.getters.getCurrency,
      memberships: Memberships.getters.getAllMemberships,
      showArchived: Memberships.getters.getShowArchived
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    Memberships.actions.fetch({locationId});
    Location.actions.loadUserLocations(uid);
    Settings.actions.loadSettings({locationId});
    Categories.actions.loadAll(locationId);
  },

  categoryChecked(id) {
    console.log('#categoryChecked(', id, ')');
    console.log('mb=', this.membership().toJS());
  },

  membership() {
    const id = parseInt(this.props.params.membershipId);

    if (this.state.memberships) {
      return this.state.memberships.find(m => m.get('Id') === id);
    } else {
      return undefined;
    }
  },

  render() {
    const mb = this.membership();
    var machine = {};

    if (this.state.categories) {
      console.log('categories=', this.state.categories.toJS());
    }

    if (!mb) {
      return <LoaderLocal/>;
    }

    return (
      <div className="container-fluid">
        <h2>Edit Membership "{mb.get('Title')}"</h2>
        <h3>Basic settings</h3>

        <hr/>

        <div className="row">

          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-title">Membership Title</label>
              <input type="text" id="membership-title" className="form-control"
                     placeholder="Membership title" value={mb.get('Title')}/>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-shortname">Shortname</label>
              <input type="text" id="membership-shortname" className="form-control"
                     placeholder="Shortname" value={mb.get('ShortName')}/>
            </div>
          </div>
          
          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-duration">Duration (months)</label>
              <input type="text" id="membership-duration" className="form-control"
                     placeholder="Duration" value={mb.get('DurationMonths')}/>
            </div>
          </div>
          
          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-price">Monthly Price</label>
              <div className="input-group">
                <input type="text" id="membership-price" className="form-control"
                       placeholder="Monthly Price" value={mb.get('MonthlyPrice')}/>
                <div className="input-group-addon">
                  {this.state.currency || '€'}
                </div>
              </div>
            </div>
          </div>

        </div>

        <div className="row">
          <div className="col-sm-6">
            <div className="form-group">
              <label>Categories Affected by Membership</label>
              <div className="row">

                {this.state.categories
                  ? (this.state.categories.map(c => {
                    return (
                      <div className="col-sm-6" key={c.get('Id')}>
                        <div className="checkbox-inline">
                          <label>
                            <input type="checkbox"
                                   value={this.categoryChecked(c.get('Id'))}
                                   ng-model="machine['Checked']"/>
                            {c.get('Name')}
                          </label>
                        </div>
                      </div>
                    );
                  })) : null}
              </div>
            </div>
          </div>
          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-price-deduction">Machine price deduction</label>
              <div className="input-group">
                <input type="text" 
                  id="membership-price-deduction" 
                  className="form-control"
                  placeholder="Percentage" 
                  value={mb.get('MachinePriceDeduction')}/>
                <div className="input-group-addon">
                  <b>%</b>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="row">

          <div className="col-sm-3">
            <div className="form-group">
              <label>Automatically extend membership</label>
              <div className="checkbox-inline">
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

          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-extend-duration">After the end date membership extends</label>
              <div className="input-group">
                <input type="text"
                  id="membership-extend-duration"
                  className="form-control"
                  value={mb.get('AutoExtendDurationMonths')}
                  ng-disabled="!membership.AutoExtend"/>
                <div className="input-group-addon">-monthly</div>
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
