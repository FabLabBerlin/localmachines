import {hashHistory} from 'react-router';
var React = require('react');


var Machine = React.createClass({
  click() {
    hashHistory.push();
  },

  imgUrl() {
    if (this.props.machine.get('Image')) {
      return '/files/' + this.props.machine.get('Image');
    } else {
      return '/machines/img/img-machine-placeholder.svg';
    }
  },

  render() {
    const m = this.props.machine;
    const style = {
      backgroundImage: 'url(' + this.imgUrl() + ')'
    };

    return (
      <a className="ms-machine" href={'/machines/#/machines/' + m.get('Id')}>
        <div className="ms-machine-label">
          <div className="ms-machine-name">
            {m.get('Name')}
          </div>
          <div className="ms-machine-brand">
            {m.get('Brand')}
          </div>
        </div>
        <div className="ms-machine-icon" style={style}>
        </div>
      </a>
    );
  }
});

export default Machine;
