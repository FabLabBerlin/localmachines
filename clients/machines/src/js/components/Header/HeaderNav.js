var Bottom = require('./Bottom');
var getters = require('../../getters');
var LoginActions = require('../../actions/LoginActions');
var {Navigation} = require('react-router');
var React = require('react');
var reactor = require('../../reactor');
var Top = require('./Top');


var HeaderNav = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isLogged: getters.getIsLogged
    };
  },

  render() {
    if (this.state.isLogged) {
      return (
        <div className="nav">
          <Top/>
          <Bottom location={this.props.location}/>
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default HeaderNav;
