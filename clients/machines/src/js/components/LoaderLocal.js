var React = require('react');


var LoaderLocal = React.createClass({
  render() {
    return (
      <div className="loader-local">
        <div className="spinner">
          <i className="fa fa-cog fa-spin"></i>
        </div>
      </div>
    );
  }
});

export default LoaderLocal;
