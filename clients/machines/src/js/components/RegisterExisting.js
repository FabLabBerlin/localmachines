var $ = require('jquery');
var getters = require('../getters');
var LoaderLocal = require('./LoaderLocal');
var Location = require('../modules/Location');
var LoginActions = require('../actions/LoginActions');
var React = require('react');
var reactor = require('../reactor');
var toastr = require('toastr');


var RegisterExisting = React.createClass({
  mixins: [ reactor.ReactMixin ],

  cancel() {
    LoginActions.logout(this.context.router);
  },

  getDataBindings() {
    return {
      location: Location.getters.getLocation,
      locationTermsUrl: Location.getters.getLocationTermsUrl,
      userId: getters.getUid
    };
  },

  componentWillMount() {
    Location.actions.loadLocations();
    var locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    Location.actions.loadTermsUrl(locationId);
  },

  render() {
    if (this.state.location) {
      return (
        <div className="container">
          <h3>Oops...</h3>
          <p>
            You are already registered on easylab.io but not a member
            of {this.state.location.get('Title')} yet. Do you want to use your
            existing Easy Lab account to use their serivces? We will give
            them access to some of your user Data that they need to
            provide you with a good service.
          </p>
          <p>
            <label>
              <input ref="acceptTerms" type="checkbox"/>
              &nbsp;
              Ja, Makea Industries GmbH darf meine Nutzerdaten
              an {this.state.location.get('Title')} weitergeben und ich habe die
              die <a href={this.state.locationTermsUrl}>AGB {this.state.location.get('Title')}</a> gelesen
              und stimme ihnen uneingeschränkt zu.
            </label>
          </p>
          <hr/>
          <div className="pull-right">
            <button className="btn btn-info btn-lg wizard-button"
                    onClick={this.cancel}>
              Cancel
            </button>
            <button className="btn btn-primary btn-lg wizard-button"
                    onClick={this.submit}>
              Continue
            </button>
          </div>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  },

  submit() {
    var acceptTerms = this.refs.acceptTerms.checked;
    if (acceptTerms) {
      var locationId = this.state.location.get('Id');
      var userId = this.state.userId;
      var router = this.context.router;

      Location.actions.addLocation({locationId, userId, router});
    } else {
      toastr.error('Accepting the AGBs is mandatory to add this location');
    }
  }
});

export default RegisterExisting;
