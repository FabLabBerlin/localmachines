var $ = require('jquery');
var getters = require('../../getters');
var MachineList = require('../MachinePage/MachineList');
var LoginStore = require('../../stores/LoginStore');
var MachineStore = require('../../stores/MachineStore');
var MachineActions = require('../../actions/MachineActions');
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
var threeDPrinterPage = React.createClass({

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
    MachineActions.apiGetUserMachines(uid);
  },

  /*
   * Initial State
   * fetch data = require(MachineStore
   */
  getDataBindings() {
    return {
      userInfo: getters.getUserInfo,
      machineInfo: getters.getMachineInfo,
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
   * Clear state while logout
   */
  clearState() {
    MachineActions.clearState();
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
    this.clearState();
    clearInterval(this.interval);
  },

  /*
   * Call when the component is mounted in DOM
   * Synchronize invent = require(stores
   * Activate a polling (1,5s)
   */
  componentDidMount() {

    this.nfcOnDidMount();
    MachineStore.onChangeActivation = this.onChangeActivation;
    LoginStore.onChangeLogout = this.onChangeLogout;
    MachineStore.onChangeLogin = this.onChangeLogin;
    this.interval = setInterval(this.update, 1500);


  },

  goToCategories(event) {
    event.preventDefault();
    window.location = '/machines/#/categories';
  },
  /*
   * Render the user name
   * MachinList
   * exit button
   */
  render() {
    var imageUrl;
    if (this.props.info && this.props.info.Image) {
      imageUrl = '/files/' + this.props.info.Image;
    } else {
      imageUrl = '/machines/img/img-machine-placeholder.svg';
    }
    var machineInfo = reactor.evaluateToJS(getters.getMachineInfo);
    if (this.state.activationInfo) {
      return (
        <div className="container-fluid">
          <div className="row">

            <div className="logged-user-name">
              <div className="text-center ng-binding">
                <i className="fa fa-user-secret"></i>&nbsp;
                {this.state.userInfo.get('FirstName')} {this.state.userInfo.get('LastName')}
              </div>
            </div>
            <div className="machine-action-info">
              <h3 className="categories-header"><img className="machine-image" src={imageUrl}/>3D Printers</h3> 
              
            </div>
            <MachineList
              user={this.state.userInfo}
              info={machineInfo}
              activation={this.state.activationInfo}
              category="3DPrinter"
            />
            <div className="container-fluid">
              <button
                onClick={this.goToCategories}
                className="btn btn-lg btn-block btn-primary btn-logout-bottom">
                Categories
              </button>
            </div>
            <div className="container-fluid">
              <button
                onClick={this.handleLogout}
                className="btn btn-lg btn-block btn-danger btn-logout-bottom">
                Sign out
              </button>
            </div>
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
    } else {
      return <div/>;
    }
  },

  /*
   * update
   *
   * Need polling for activation status and maintenance status
   */
  update() {
    MachineActions.pollActivations();
    MachineActions.pollMachines();
  }
});

export default threeDPrinterPage;