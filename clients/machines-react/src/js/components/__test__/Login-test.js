jest.mock('jquery');
jest.mock('../../stores/LoginStore.js');
jest.dontMock('../Login.js');
jest.dontMock('../../actions/LoginActions.js');


describe('Login', function() {
  it('renders Log In and Sign Up', function() {
    var React = require('react');
    var Login = React.createFactory(require('../Login'));
    var login = new Login({});
    var s = React.renderToString(login);
    expect(s).toContain('Log In');
    expect(s).toContain('Sign up');
  });

  it('does something when clicking Log In', function() {
    var React = require('react/addons');
    var TestUtils = React.addons.TestUtils;
    var Login = require('../Login');
    var LoginStore = require('../../stores/LoginStore.js');
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
    expect(LoginStore.apiPostLogin).toBeCalledWith({
      username: 'joe',
      password: '123456'
    });
  });
});
