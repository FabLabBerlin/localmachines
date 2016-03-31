jest.dontMock('jquery');
jest.mock('../../../actions/LoginActions.js');
jest.dontMock('../Login.js');
jest.dontMock('../LoginChooser.js');
jest.dontMock('../LoginNfc.js');

describe('LoginChooser', function() {
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
    expect(s).toContain('class="nfc-login-icon"');
    expect(s).toContain('class="nfc-login-info-text"');
    expect(s).toContain('Use your NFC card to log in');
    window.libnfc = null;
  });

});
