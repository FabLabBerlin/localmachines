var LocationGetters = require('../../modules/Location/getters');
var React = require('react');
var reactor = require('../../reactor');


var Top = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation
    }
  },

  render() {
    return (
      <div className="row">
        <div className="col-xs-12">
          <div className="nav-top">
            <a className="nav-logo" 
               href="/machines/#/machine">
              <img src="/machines/assets/img/logo-easylab.svg" className="brand-image hidden-xs"/>
              <img src="/machines/assets/img/logo-small.svg" className="brand-image visible-xs-block"/>
            </a>
            <div className="nav-title">
              {this.state.location ? (
                <span className="nav-title-lab">
                  {this.state.location.get('Title')}
                </span>
              ) : null}
              <span className="nav-title-easylab">EASY LAB</span>
            </div>
          </div>
        </div>
      </div>
    );
  }
});

export default Top;
