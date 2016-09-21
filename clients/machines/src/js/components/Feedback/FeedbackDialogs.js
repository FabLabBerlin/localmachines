var $ = require('jquery');
var FeedbackActions = require('../../actions/FeedbackActions');
var getters = require('../../getters');
var LoginActions = require('../../actions/LoginActions');
var Machines = require('../../modules/Machines');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

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
          click($vexContent, event) {
            $vexContent.data().vex.value = 'good';
            return vex.close($vexContent.data().vex.id);
          }
        },
        NEUTRAL: {
          text: 'Neutral',
          type: 'button',
          className: 'vex-dialog-button-primary',
          click($vexContent, event) {
            $vexContent.data().vex.value = 'neutral';
            return vex.close($vexContent.data().vex.id);
          }
        },
        BAD: {
          text: 'Bad',
          type: 'button',
          className: 'vex-dialog-button-primary',
          click($vexContent, event) {
            $vexContent.data().vex.value = 'bad';
            return vex.close($vexContent.data().vex.id);
          }
        }
      },
      message: 'Satisfied with the outcome?',
      afterOpen($vexContent) {
        $vexContent.find('button[type="submit"]').hide();
      },
      callback(satisfaction) {
        if (satisfaction) {
          $('.vex').remove();
          $('body').removeClass('vex-open');
          console.log('satisfaction:', satisfaction);
          FeedbackActions.reportSatisfaction({ activationId, satisfaction });
        }
      }
    });
  },

  machineIssue(machineId) {
    LoginActions.keepAlive();
    const machinesById = reactor.evaluateToJS(Machines.getters.getMachinesById);
    const machine = machinesById[machineId] || {};
    VexDialog.buttons.YES.text = 'Please fix it';
    VexDialog.buttons.NO.text = 'Nevermind';

    VexDialog.prompt({
      message: 'What happened?',
      placeholder: 'I saw...',
      callback: function(text) {
        if (text) {
          FeedbackActions.reportMachineBroken({ machineId, text });
        } else if (text !== false) {
          toastr.error('Please give us some information.');
        }
      }
    });
    /*VexDialog.confirm({
      message: 'Do you really want to report machine <b>' +
        machine.Name + '</b> as broken?',
      callback(confirmed) {
        if (confirmed) {
          FeedbackActions.reportMachineBroken({ machineId });
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }
    });*/
  }
};
