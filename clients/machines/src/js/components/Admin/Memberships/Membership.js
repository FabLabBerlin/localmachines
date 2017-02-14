var _ = require('lodash');
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

  handleChange(key, typeConverter, e) {
    const value = typeConverter(e.target.value);
    console.log('value=', value);
    Memberships.actions.setMembershipField(
      this.membership().get('Id'),
      key, 
      _.isNaN(value) ? e.target.value : value
    );
  },

  handleSave() {
    Memberships.actions.save(this.membership().get('Id'));
  },

  categoryChecked(id) {
    if (!this.membership().get('AffectedCategories')) {
      return false;
    }

    const ids = JSON.parse(this.membership().get('AffectedCategories'));

    return _.includes(ids, id);
  },

  setCategoryChecked(id, yes) {
    Memberships.actions.setMembershipCategory(this.membership().get('Id'), id, yes);
  },

  setArchive(yes) {
    Memberships.actions.setMembershipArchive(this.membership().get('Id'), yes);
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
              <input type="text"
                     id="membership-title"
                     className="form-control"
                     onChange={this.handleChange.bind(this, 'Title', String)}
                     placeholder="Membership title"
                     value={mb.get('Title')}/>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-shortname">Shortname</label>
              <input type="text" 
                     id="membership-shortname" 
                     className="form-control"
                     onChange={this.handleChange.bind(this, 'ShortName', String)}
                     placeholder="Shortname" 
                     value={mb.get('ShortName')}/>
            </div>
          </div>
          
          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-duration">Duration (months)</label>
              <input type="text" 
                     id="membership-duration" 
                     className="form-control"
                     onChange={this.handleChange.bind(this, 'DurationMonths', parseInt)}
                     placeholder="Duration" 
                     value={mb.get('DurationMonths')}/>
            </div>
          </div>
          
          <div className="col-sm-3">
            <div className="form-group">
              <label htmlFor="membership-price">Monthly Price</label>
              <div className="input-group">
                <input type="text" 
                       id="membership-price" 
                       className="form-control"
                       onChange={this.handleChange.bind(this, 'MonthlyPrice', parseFloat)}
                       placeholder="Monthly Price" 
                       value={mb.get('MonthlyPrice')}/>
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
                    const checked = this.categoryChecked(c.get('Id'));

                    return (
                      <div className="col-sm-6" key={c.get('Id')}>
                        <div className="checkbox-inline">
                          <label>
                            <input type="checkbox"
                                   checked={checked}
                                   onChange={this.setCategoryChecked.bind(this, c.get('Id'), !checked)}/>
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
                  onChange={this.handleChange.bind(this, 'MachinePriceDeduction', parseFloat)}
                  placeholder="Percentage" 
                  value={mb.get('MachinePriceDeduction')}/>
                <div className="input-group-addon">
                  <b>%</b>
                </div>
              </div>
            </div>
          </div>
        </div>

        <hr/>

  <div className="pull-right">

    {(this.membership() && this.membership().get('Archived')) ? (
      <button className="btn btn-danger btn-lg"
              onClick={this.setArchive.bind(this, false)}>
        <i className="fa fa-archive"></i>&nbsp;Unarchive
      </button>
    ) : (
      <button className="btn btn-danger btn-lg"
              onClick={this.setArchive.bind(this, true)}>
        <i className="fa fa-archive"></i>&nbsp;Archive
      </button>
    )}

    <button className="btn btn-primary btn-lg"
            onClick={this.handleSave}>
      <i className="fa fa-save"></i>&nbsp;Save
    </button>

  </div>

      </div>
    );
  }
});

export default MembershipPage;
