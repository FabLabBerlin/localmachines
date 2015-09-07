var FeedbackActions = require('../actions/FeedbackActions');
var getters = require('../getters');
var {Navigation} = require('react-router');
var NfcLogoutMixin = require('./NfcLogoutMixin');
var React = require('react');
var reactor = require('../reactor');
var UserActions = require('../actions/UserActions');


var FeedbackPage = React.createClass({
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

  getDataBindings() {
    return {
      subject: getters.getFeedbackSubject,
      subjectDropdown: getters.getFeedbackSubjectDropdown,
      subjectOtherText: getters.getFeedbackSubjectOtherText,
      message: getters.getFeedbackMessage
    };
  },

  componentDidMount() {
    this.nfcOnDidMount();
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.getUserInfoFromServer(uid);
  },

  componentWillUnmount() {
    this.nfcOnWillUnmount();
  },

  handleChange(e) {
    var key = event.target.id;
    var value = event.target.value;
    FeedbackActions.setFeedbackProperty({ key, value });
  },

  handleSubmit() {
    FeedbackActions.submit();
  },

  render() {
    return (
      <div>
        <div>
          <label>
            Subject:
            <select id="subject-dropdown"
                    onChange={this.handleChange}
                    value={this.state.subjectDropdown}>
              <option>Billing</option>
              <option>Technical</option>
              <option>Other</option>
            </select>
          </label>
        </div>
        {this.state.subjectDropdown === 'Other' ?
          (
            <div>
              <label>
                Other text:
                <input id="subject-other-text"
                       onChange={this.handleChange}
                       value={this.state.subjectOtherText}/>
              </label>
            </div>
          ) : null
        }
        <div>
          <label>
            Message:
            <textarea id="message"
                      onChange={this.handleChange}
                      value={this.state.message}
            />
          </label>
        </div>
        <div>
          <button onClick={this.handleSubmit}>Submit</button>
        </div>
      </div>
    );
  }
});

export default FeedbackPage;
