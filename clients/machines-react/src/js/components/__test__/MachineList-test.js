var React = require('react');
var MachineList = React.createFactory(require('../MachineList'));


function machines() {
  return [
    {
      Id: 1,
      Name: 'Printer5000'
    },
    {
      Id: 2,
      Name: 'Form1'
    }
  ];
}


describe('MachineList', function() {
  describe('render', function() {
    it('says access denied when no machines are passed as parameter', function() {
      var ml = new MachineList({
        info: []
      });
      var s = React.renderToString(ml);
      expect(s).toContain('You do not have access');
    });

    it('renders the machines that are passed as parameter', function() {
      var ml = new MachineList({
        info: machines(),
        user: {
          Role: 'member'
        }
      });
      var s = React.renderToString(ml);
      expect(s).toContain('Printer5000');
      expect(s).toContain('Form1');
    });

    it('shows the force on/off switch if and only if the user is admin', function() {
      [false, true].forEach(function(isAdmin) {
        var ml = new MachineList({
          info: machines(),
          user: {
            Role: isAdmin ? 'admin' : 'member'
          }
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
