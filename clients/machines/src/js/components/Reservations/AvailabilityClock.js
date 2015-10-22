var $ = require('jquery');
var moment = require('moment');
var React = require('react');


var HOURS = 12;
var SLOTS_PER_HOUR = 2;
var SLOTS = HOURS * SLOTS_PER_HOUR;


var COLOR_CURRENT_TIME = '#ed9393';
var COLOR_AVAILABLE = '#b8d879';


var AvailabilityClock = React.createClass({
  componentDidMount() {
    this.drawCanvas();
  },

  componentDidUpdate() {
    //this.drawCanvas();
  },

  componentWillUmount() {
    var n = this.refs.canvasContainer.getDOMNode();
    $(n).find('canvas').remove();
  },

  drawCanvas() {
    var n = this.refs.canvasContainer.getDOMNode();
    $(n).append('<canvas width="80" height="80"></canvas>');
    var canvas = $(n).find('canvas')[0];
    var context = canvas.getContext('2d');

    for (var i = 0; i < SLOTS; i++) {
      if ((SLOTS_PER_HOUR * moment().hour()) % SLOTS === i) {
        this.drawSlot(context, i, COLOR_CURRENT_TIME);
      } else {
        this.drawSlot(context, i, COLOR_AVAILABLE);
      }
    }
  },

  drawSlot(context, i, color) {
    var startAngle = 360 / SLOTS * i;
    var endAngle = startAngle + 360 / SLOTS;
    startAngle -= 90;
    startAngle = (startAngle + 360) % 360;
    startAngle *= Math.PI / 180;
    endAngle *= Math.PI / 180;

    context.beginPath();
    context.moveTo(40, 40);
    context.arc(40, 40, 30, startAngle, endAngle, false);
    context.lineTo(40, 40);
    context.closePath();
    context.fillStyle = color;
    context.fill();
  },

  render() {
    return <div ref="canvasContainer"/>;
  }
});

export default AvailabilityClock;
