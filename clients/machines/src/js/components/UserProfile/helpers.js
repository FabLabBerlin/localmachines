
/*
 * Date utility function
 *
 * expecting date types moment.js
 */
function endDate(startDate, duration) {
  return startDate.add(duration, 'd');
}

function formatDate(d) {
  return d.format('DD. MMM YYYY');
}

export default {endDate, formatDate};
