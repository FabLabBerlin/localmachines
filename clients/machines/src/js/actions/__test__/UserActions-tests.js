jest.dontMock('../ApiActions');
jest.dontMock('../../actionTypes');
jest.dontMock('../UserActions');
jest.dontMock('nuclear-js');
jest.mock('jquery');
jest.mock('../../reactor');

var $ = require('jquery');
var actionTypes = require('../../actionTypes');
var reactor = require('../../reactor');
var UserActions = require('../UserActions');

const uidTest = 5;

const loginInfoTest = {
  username: 'test',
  password: 'test'
};

const userUpdateState = {
  FirstName: 'user',
  LastName: 'Update',
  Username: 'State',
  Email: 'user@example.com',
  InvoiceAddr: 0,
  ShipAddr: 0
};

describe('UserActions', function() {
  describe('getUserInfoFromServer', function() {
    it('should GET /api/users/:uid', function() {
      UserActions.getUserInfoFromServer(123);
      expect($.ajax).toBeCalledWith({
        url: '/api/users/123',
        dataType: 'json',
        type: 'GET',
        cache: false,
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  it('test getInfoBillFromServer', function() {
    UserActions.getInfoBillFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5/bill',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  it('test getMembershipFromServer', function() {
    UserActions.getMembershipFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5/memberships',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  it('test updatePassword', function() {
    UserActions.updatePassword(uidTest, 'passwordTest');
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5/password',
      dataType: 'json',
      type: 'POST',
      data: {
        password: 'passwordTest'
      },
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  it('test submitState', function() {
    UserActions.submitState(uidTest, userUpdateState);
    expect($.ajax).toBeCalledWith({
      headers: {'Content-Type': 'application/json'},
      url: '/api/users/5',
      type: 'PUT',
      data: JSON.stringify({
        User: userUpdateState
      }),
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

});