var $ = require('jquery');
var getters = require('../../getters');
var LoginStore = require('../../stores/LoginStore');
var NfcLogoutMixin = require('../Login/NfcLogoutMixin');
var LoginActions = require('../../actions/LoginActions');
var Navigation = require('react-router').Navigation;
var React = require('react');
var reactor = require('../../reactor');
var ScrollNav = require('../ScrollNav');
var toastr = require('../../toastr');
var UserActions = require('../../actions/UserActions');

/*
 * MachinePage:
 * Root component
 * Fetch the information = require(the store
 * Give it to its children to display the interface
 * TODO: reorganize and documente some function
 */
var CategoriesPage = React.createClass({

  /*
   * Enable some React router function as:
   *  ReplaceWith
   */
  mixins: [ Navigation, reactor.ReactMixin, NfcLogoutMixin ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      const isLogged = reactor.evaluateToJS(getters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  /*
   * Start fetching the data
   * before the component is mounted
   */
  componentWillMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.getUserInfoFromServer(uid);
  },

  /*
   * Initial State
   * fetch data = require(MachineStore
   */
  getDataBindings() {
    return {
      userInfo: getters.getUserInfo,
      activationInfo: getters.getActivationInfo,
      isLoading: getters.getIsLoading
    };
  },

  /*
   * Create a table of the Id = require(an array
   * Used in shouldComponentUpdate to know get the id = require(previous state and next one
   */
  createCompareTable(state) {
    let table = [];
    for(let i in state) {
      table.push(state[i].Id);
    }
    return table;
  },

  /*
   * Look if the activations have a name
   * if they all have one, return true
   */
  hasNameInto(activation) {
    for(let i in activation ) {
      if(!activation.FirstName) {
        return false;
      }
    }
    return true;
  },


  /*
   * Logout with the exit button
   */
  handleLogout() {
    LoginActions.logout(this.context.router);
  },


  /*
   * Destructor
   * Stop the polling
   */
  componentWillUnmount() {
    this.nfcOnWillUnmount();
    clearInterval(this.interval);
  },

  /*
   * Call when the component is mounted in DOM
   * Synchronize invent = require(stores
   * Activate a polling (1,5s)
   */
  componentDidMount() {

    this.nfcOnDidMount();
    LoginStore.onChangeLogout = this.onChangeLogout;
    this.interval = setInterval(this.update, 1500);


  },

  goTo3Dprinters(event) {
    event.preventDefault();
    window.location = '/machines/#/3dprinter';
  },

  goToLasercutter(event) {
    event.preventDefault();
    window.location = '/machines/#/lasercutter';
  },
  /*
   * Render the user name
   * Categories
   * exit button
   */
  render() {
  
      return (
        <div className="containter-fluid">
          <div className="logged-user-name">
            <div className="text-center ng-binding">
              <i className="fa fa-user-secret"></i>&nbsp;
              {this.state.userInfo.get('FirstName')} {this.state.userInfo.get('LastName')}
            </div>
          </div>
          <div className="container-fluid">
          	<div className="row">
          		<button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4" href="#"
              onClick={this.goTo3Dprinters}>3D PRINTING
              </button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4" href="#"
              onClick={this.goToLasercutter}>LASER CUTTING</button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4">ELECTRONICS</button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4">VINYL & TEXTILES</button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4">WOOD WORK</button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4">METAL WORK</button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4">bads</button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4">8</button>
              <button className="btn-primary btn-cat col-xs-6 col-sm-4 col-md-4 col-lg-4">9</button>
          	</div>
          </div>
          
          <div className="container-fluid">
            <button
              onClick={this.handleLogout}
              className="btn btn-lg btn-block btn-danger btn-logout-bottom">
              Sign out
            </button>
          </div>
          <ScrollNav/>
          {
            this.state.isLoading ?
              (
              <div id="loader-global">
                <div className="spinner">
                  <i className="fa fa-cog fa-spin"></i>
                </div>
              </div>
            )
            : ''
          }
        </div>
      );
  },
 
});

export default CategoriesPage;