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
    window.location = '/machines/#/categories/3DPrinter';
  },

  goToLasercutter(event) {
    event.preventDefault();
    window.location = '/machines/#/categories/Lasercutter';
  },

  goToWoodWork(event) {
    event.preventDefault();
    window.location = '/machines/#/categories/WoodWork';
  },

  goToMetalWork(event) {
    event.preventDefault();
    window.location = '/machines/#/categories/MetalWork';
  },

  goToVinylnTextiles(event) {
    event.preventDefault();
    window.location = '/machines/#/categories/VinylnTextiles';
  },

  goToElectronics(event) {
    event.preventDefault();
    window.location = '/machines/#/categories/Electronics';
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
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
          		  <div className="category-pad row">
                  <button className="button-category" href="#"
                  onClick={this.goTo3Dprinters}>
                    <div>
                      <img className="machine-image" src='../../../../../../files/machine-4.svg'/>
                    </div>
                    <div>
                      3D Printing
                    </div>
                  </button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row">
                  <button className="button-category" href="#"
                  onClick={this.goToLasercutter}>
                    <div>
                      <img className="machine-image" src='../../../../../../files/machine-3.svg'/>
                    </div>
                    <div>
                      Laser Cutting
                    </div>
                  </button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row">                 
                  <button className="button-category" href="#"
                  onClick={this.goToElectronics}>
                    <div>
                      <img className="machine-image" src='../../../../../../files/machine-9.svg'/>
                    </div>
                    <div>
                      Electronics
                    </div>                  
                  </button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row">
                  <button className="button-category" href="#"
                  onClick={this.goToVinylnTextiles}>
                    <div>
                      <img className="machine-image" src='../../../../../../files/machine-14.svg'/>
                    </div>
                    <div>
                      Vinyl & Textiles
                    </div>                  
                  </button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row">
                  <button className="button-category" href="#"
                  onClick={this.goToWoodWork}>
                    <div>
                      <img className="machine-image" src='../../../../../../files/machine-11.svg'/>
                    </div>
                    <div>
                      Wood Work
                    </div>
                  </button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row">
                  <button className="button-category" href="#"
                  onClick={this.goToMetalWork}>METAL WORK</button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row">
                  <button className="button-category">bads</button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row"> 
                  <button className="button-category">8</button>
                </div>
              </div>
              <div className="col-xs-6 col-sm-4 col-md-4 col-lg-4">
                <div className="category-pad row">
                  <button className="button-category">9</button>
                </div>
              </div>
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