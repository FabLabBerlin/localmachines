angular.module('fabsmithFilters', [])

.filter('timeleft', function() {
  return function(inputSeconds) {

  	var hours = Math.floor(inputSeconds / 3600);
		inputSeconds = inputSeconds - hours * 3600;
  	var minutes = Math.floor(inputSeconds / 60);
  	var seconds = inputSeconds - minutes * 60;

  	var timeStr = '';

  	if (hours > 0) {
  		var hoursStr = '';
  		if (hours <= 9) {
  			hoursStr = '0';
  		}
  		hoursStr += hours.toString();
  		timeStr += hoursStr + ':';
  	}

  	if (minutes > 0) {
  		var minutesStr = '';
  		if (minutes <= 9) {
  			minutesStr = '0';
  		}
  		minutesStr += minutes.toString();
  		timeStr += minutesStr + ':';
  	}

  	var secondsStr = '';
  	if (seconds <= 9) {
  		secondsStr = '0';
  	}
  	secondsStr += seconds.toString();
  	timeStr += secondsStr;

    return timeStr;
  };
});