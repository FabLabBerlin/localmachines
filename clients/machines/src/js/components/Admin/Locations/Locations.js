var React = require('react');
var reactor = require('../../../reactor');


var Locations = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      //isLogged: getters.getIsLogged
    };
  },

  render() {
    return (
      <div className="container">
        <h1>Location List</h1>
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Id</th>
              <th>Title</th>
              <th>First Name</th>
              <th>Last Name</th>
              <th>E-Mail</th>
              <th>Jabber ID</th>
            </tr>
          </thead>
          <tbody>
            <tr>
            </tr>
          </tbody>
        </table>
      </div>
    );
  }
});

export default Locations;
