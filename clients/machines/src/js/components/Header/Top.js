var React = require('react');


var Top = React.createClass({
  render() {
    return (
      <div className="row">
        <div className="col-xs-12">
          <div className="navbar-brand">
            <a 
              className="brand-link" 
              href="/machines/#/machine">
              <img src="/machines/assets/img/logo-easylab.svg" className="brand-image hidden-xs"/>
              <img src="/machines/assets/img/logo-small.svg" className="brand-image visible-xs-block"/>
            </a>
          </div>
        </div>
      </div>
    );
  }
});

export default Top;
