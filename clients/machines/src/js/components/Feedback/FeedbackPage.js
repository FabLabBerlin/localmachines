var FeedbackActions = require('../../actions/FeedbackActions');
var getters = require('../../getters');
var LocationActions = require('../../actions/LocationActions');
var React = require('react');
var reactor = require('../../reactor');
var UserActions = require('../../actions/UserActions');


var FeedbackPage = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      subject: getters.getFeedbackSubject,
      subjectDropdown: getters.getFeedbackSubjectDropdown,
      subjectOtherText: getters.getFeedbackSubjectOtherText,
      message: getters.getFeedbackMessage
    };
  },

  componentDidMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    LocationActions.loadUserLocations(uid);
  },

  handleChange(event) {
    var key = event.target.id;
    var value = event.target.value;
    FeedbackActions.setFeedbackProperty({ key, value });
  },

  handleSubmit() {
    FeedbackActions.submit();
  },

  render() {
    return (
      <div className="container">
        <h3>Your Feedback</h3>
        <div className="form-group">
          <label htmlFor="subject-dropdown">
            Subject
          </label>
          <select id="subject-dropdown"
                  className="form-control"
                  onChange={this.handleChange}
                  value={this.state.subjectDropdown}>
            <option>Billing</option>
            <option>Technical</option>
            <option>Other</option>
          </select>
        </div>
        {this.state.subjectDropdown === 'Other' ?
          (
            <div className="form-group">
              <div>
                <label htmlFor="subject-other-text">
                  Other text
                </label>
                <input id="subject-other-text"
                       className="form-control"
                       onChange={this.handleChange}
                       value={this.state.subjectOtherText}/>
              </div>
            </div>
          ) : null
        }
        <div className="form-group">
          <label htmlFor="message">
            Message
          </label>
          <textarea id="message"
                    className="form-control"
                    onChange={this.handleChange}
                    value={this.state.message}
          />
        </div>
        <div>
          <button className="btn btn-primary btn-lg"
                  onClick={this.handleSubmit}>Submit</button>
        </div>
      </div>
    );
  }
});

export default FeedbackPage;
