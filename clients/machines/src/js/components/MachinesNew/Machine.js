import {hashHistory} from 'react-router';
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');


const AVAILABLE = 'available';
const LOCKED = 'locked';
const MAINTENANCE = 'maintenance';
const OCCUPIED = 'occupied';
const RESERVED = 'reserved';
const RUNNING = 'running';
const STAFF = 'staff';


var Machine = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      user: getters.getUser
    };
  },

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
    if (m.get('activation')) {
      console.log('machine with activation:', m.toJS());
    }

    return (
      <a className={'ms-machine ' + this.statusClass()}
         href={'/machines/#/machines/' + m.get('Id')}>
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
  },

  status() {
    const m = this.props.machine;
    const a = m.get('activation');

    if (m.get('UnderMaintenance')) {
      return MAINTENANCE;
    } else {
      if (a) {
        console.log('machine ' + m.get('Name'));
        console.log('a=', a);
        console.log('a.UserId=', a.get('UserId'));
        console.log('userId=', this.state.user.get('Id'));
        console.log('this.state.user=', this.state.user);
        if (a.get('UserId') === this.state.user.get('Id')) {
          return RUNNING;
        } else {
          return OCCUPIED;
        }
      } else {
        return AVAILABLE;
      }
    }
  },

  statusClass() {
    return 'ms-' + this.status();
  }
});

export default Machine;
