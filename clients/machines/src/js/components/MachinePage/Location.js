var _ = require('lodash');
var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var LocationActions = require('../../actions/LocationActions');
var React = require('react');
var reactor = require('../../reactor');


var Location = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      locations: getters.getLocations
    };
  },

  componentWillMount() {
    LocationActions.loadLocations();
  },

  handleChange(event) {
    var id = parseInt(event.target.value);
    LocationActions.setLocationId(id);
  },

  render() {
    if (this.state.locations) {
      return (
        <div className="container-fluid location-picker">
        <div className="form-horizontal">
          <div className="form-group">
          <label 
            htmlFor="location"
            className="col-sm-2 control-label">Location</label>
          <div className="col-sm-10">
          <select
            id="location" 
            className="form-control" 
            onChange={this.handleChange}>
            {_.map(this.state.locations, (location, i) => {
              if (location.Approved) {
                return (
                  <option key={i} value={location.Id}>
                    {location.Title}
                  </option>
                );
              }
            })}
          </select>
          </div>
          </div>
        </div>
        </div>
      );
    } else {
      return <div></div>;
    }
  }

});

export default Location;
