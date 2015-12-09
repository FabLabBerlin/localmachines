jest.dontMock('../ForceSwitch');
jest.dontMock('../FreeMachine');
jest.dontMock('../MachineChooser');
jest.dontMock('../MachineList');
jest.dontMock('lodash');
jest.dontMock('nuclear-js');
jest.dontMock('react');

var React = require('react');
var MachineList = React.createFactory(require('../MachineList'));
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


function machines() {
  return [
    {
      Id: 1,
      Name: 'Printer5000',
      Visible: true
    },
    {
      Id: 2,
      Name: 'Form1',
      Visible: true
    }
  ];
}


describe('MachineList', function() {
  describe('render', function() {
    it('says access denied when no machines are passed as parameter', function() {
      var ml = new MachineList({
        machines: []
      });
      var s = React.renderToString(ml);
      expect(s).toContain('You do not have access');
    });

    it('renders the machines that are passed as parameter', function() {
      var ml = new MachineList({
        machines: machines(),
        user: toImmutable({
          UserRole: 'member'
        })
      });
      var s = React.renderToString(ml);
      expect(s).toContain('Printer5000');
      expect(s).toContain('Form1');
    });

    it('shows the force on/off switch if and only if the user is admin', function() {
      [false, true].forEach(function(isAdmin) {
        var ml = new MachineList({
          machines: machines(),
          user: toImmutable({
            UserRole: isAdmin ? 'admin' : 'member'
          })
        });
        var s = React.renderToString(ml);
        expect(s).toContain('Printer5000');
        if (isAdmin) {
          expect(s).toContain('Force Switch');
        } else {
          expect(s).not.toContain('Force Switch');
        }
      });
    });
  });
});
