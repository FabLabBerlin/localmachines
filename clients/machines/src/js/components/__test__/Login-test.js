jest.mock('jquery');
jest.mock('../../actions/LoginActions.js');
jest.dontMock('../Login.js');
jest.dontMock('../LoginChooser.js');
jest.dontMock('../LoginNfc.js');

describe('LoginChooser', function() {
  it('renders Log In and Sign up in case of normal browser', function() {
    var React = require('react');
    var LoginChooser = React.createFactory(require('../LoginChooser'));
    var loginChooser = new LoginChooser({});
    var s = React.renderToString(loginChooser);
    expect(s).toContain('Log In');
    expect(s).toContain('Sign up');
  });

  it('renders NFC icon and related text in case of NFC browser', function() {
    window.libnfc = {
      debug: true,
      cardRead: { connect: function() {} },
      cardReaderError: { connect: function() {} },
      asyncScan: function() {}
    };
    var React = require('react');
    var LoginChooser = React.createFactory(require('../LoginChooser'));
    var loginChooser = new LoginChooser({});
    var s = React.renderToString(loginChooser);
    expect(s).toContain('class="nfc-login-info-icon"');
    expect(s).toContain('class="nfc-login-info-text"');
    expect(s).toContain('Use your NFC card to log in');
    window.libnfc = null;
  });

  it('does something when clicking Log In', function() {
    var React = require('react/addons');
    var TestUtils = React.addons.TestUtils;
    var Login = require('../Login');
    var LoginActions = require('../../actions/LoginActions.js');
    var login = TestUtils.renderIntoDocument(
      <Login />
    );
    var form = TestUtils.findRenderedDOMComponentWithTag(login, 'form');
    var username = TestUtils.findAllInRenderedTree(login, function(c) {
      return c.getDOMNode().getAttribute('name') === 'username';
    })[0];
    var password = TestUtils.findAllInRenderedTree(login, function(c) {
      return c.getDOMNode().getAttribute('name') === 'password';
    })[0];
    var button = TestUtils.findRenderedDOMComponentWithClass(login, 'btn-login');
    expect(button.props.children).toEqual('Log In');
    TestUtils.Simulate.change(username, {
      target: {
        name: 'username',
        value: 'joe'
      }
    });
    TestUtils.Simulate.change(password, {
      target: {
        name: 'password',
        value: '123456'
      }
    });
    TestUtils.Simulate.submit(form);
    expect(LoginActions.submitLoginForm).toBeCalledWith({
      username: 'joe',
      password: '123456'
    }, undefined);
  });

});
