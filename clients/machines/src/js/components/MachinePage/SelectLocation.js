var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var React = require('react');
var reactor = require('../../reactor');


var SelectLocation = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      locations: getters.getLocations
    };
  },

  componentWillMount() {
    MachineActions.loadLocations();
  },

  handleChange(event) {
    var id = parseInt(event.target.value);
    MachineActions.setLocationId(id);
  },

  render() {
    if (this.state.locations) {
      return (
        <div className="container-fluid">
          <select className="form-control" onChange={this.handleChange}>
            {_.map(this.state.locations, (location) => {
              if (location.Approved) {
                return (
                  <option value={location.Id}>
                    {location.Title}
                  </option>
                );
              }
            })}
          </select>
        </div>
      );
    } else {
      return <div></div>;
    }
  }

});

export default SelectLocation;
