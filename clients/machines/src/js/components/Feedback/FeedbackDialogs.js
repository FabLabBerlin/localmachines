var FeedbackActions = require('../../actions/FeedbackActions');
var getters = require('../../getters');
var reactor = require('../../reactor');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';


export default {
  checkSatisfaction(activationId) {
    VexDialog.buttons.NO.text = 'Cancel';

    VexDialog.open({
      buttons: {
        GOOD: {
          text: 'Good',
          type: 'button',
          className: 'vex-dialog-button-primary',
          click: function($vexContent, event) {
            $vexContent.data().vex.value = 'good';
            return vex.close($vexContent.data().vex.id);
          }
        },
        NEUTRAL: {
          text: 'Neutral',
          type: 'button',
          className: 'vex-dialog-button-primary',
          click: function($vexContent, event) {
            $vexContent.data().vex.value = 'neutral';
            return vex.close($vexContent.data().vex.id);
          }
        },
        BAD: {
          text: 'Bad',
          type: 'button',
          className: 'vex-dialog-button-primary',
          click: function($vexContent, event) {
            $vexContent.data().vex.value = 'bad';
            return vex.close($vexContent.data().vex.id);
          }
        }
      },
      message: 'Satisfied with the outcome?',
      afterOpen: function($vexContent) {
        $vexContent.find('button[type="submit"]').hide();
      },
      callback: function(satisfaction) {
        if (satisfaction) {
          $('.vex').remove();
          $('body').removeClass('vex-open');
          console.log('satisfaction:', satisfaction);
          FeedbackActions.reportSatisfaction({ activationId, satisfaction });
        }
      }.bind(this)
    });
  },

  machineIssue(machineId) {
    const machinesById = reactor.evaluateToJS(getters.getMachinesById);
    const machine = machinesById[machineId] || {};
    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to report machine <b>' +
        machine.Name + '</b> as broken?',
      callback: function(confirmed) {
        if (confirmed) {
          FeedbackActions.reportMachineBroken({ machineId });
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }.bind(this)
    });
  }
};
