import Bottom from './Bottom';
import BottomMachine from './BottomMachine';
import getters from '../../getters';
import LoginActions from '../../actions/LoginActions';
var {Navigation} = require('react-router');
import React from 'react';
import reactor from '../../reactor';
import Top from './Top';
import TopMachine from './TopMachine';


var HeaderNav = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isLogged: getters.getIsLogged
    };
  },

  render() {
    if (this.state.isLogged) {
      const pathMatch = this.props.location.pathname.match(/.machines.(\d+)/);
      const className = pathMatch ? 'nav-machine' : '';
      const machineId = pathMatch ? parseInt(pathMatch[1]) : undefined;

      return (
        <div className={'nav ' + className}>
          {pathMatch ?
            <TopMachine machineId={machineId}/> :
            <Top/>
          }
          {pathMatch ?
            <BottomMachine location={this.props.location}
                           machineId={machineId}/> :
            <Bottom location={this.props.location}/>
          }
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default HeaderNav;
