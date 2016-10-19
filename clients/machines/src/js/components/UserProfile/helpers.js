var moment = require('moment');
var reactor = require('../../reactor');
var SettingsGetters = require('../../modules/Settings/getters');


/*
 * VAT utility functions
 */

function addVAT(priceExclVAT) {
  var vat = reactor.evaluateToJS(SettingsGetters.getVatPercent) / 100;
  return priceExclVAT * (1 + vat);
}

function subtractVAT(priceInclVAT) {
  var vat = reactor.evaluateToJS(SettingsGetters.getVatPercent) / 100;
  return priceInclVAT / (1 + vat);
}

/*
 * toCents converts a price in euro into cents.
 *
 * It's better to do arithmetic in cents because IEEE 754
 * can do funny stuff. http://stackoverflow.com/q/588004
 */
function toCents(priceEuro) {
  return Math.round(priceEuro * 100);
}

function toEuro(priceCents) {
  return (priceCents / 100).toFixed(2);
}

/*
 * Date utility functions
 *
 * expecting date types moment.js
 */
function formatDate(d) {
  return d.format('DD. MMM YYYY');
}

function formatDuration(purchase) {
  if (purchase.Quantity) {
    var duration = purchase.Quantity;
    switch (purchase.PriceUnit) {
    case 'month':
      duration *= 60 * 60 * 24 * 30;
      break;
    case 'day':
      duration *= 60 * 60 * 24;
      break;
    case 'hour':
      duration *= 60 * 60;
      break;
    case '30 minutes':
      duration *= 60 * 30;
      break;
    case 'minute':
      duration *= 60;
      break;
    case 'second':
      break;
    default:
      console.log('unknown price unit', purchase.PriceUnit);
      return undefined;
    }

    var d = parseInt(duration.toString(), 10);
    var h = String(Math.floor(d / 3600));
    var m = String(Math.floor(d % 3600 / 60));

    if (m.length === 1) {
      m = '0' + m;
    }
    var str = h + ':' + m;

    return str;
  }
}

function formatPrice(price) {
  return (Math.round(price * 100) / 100).toFixed(2);
}

function toQuantity(p, duration) {
  var m;

  if (duration.indexOf(':') > 0) {
    m = moment.duration(duration);
  } else {
    m = moment.duration({
      hours: duration
    });
  }

  switch (p.PriceUnit) {
  case 'second':
    return m.asSeconds();
  case 'minute':
    return m.asMinutes();
  case '30 minutes':
    return m.asHours() * 2;
  case 'hour':
    return m.asHours();
  case 'day':
    return m.asDays();
  }
}

export default {
  addVAT, subtractVAT, toCents, toEuro,
  formatDate, formatDuration, formatPrice, toQuantity
};
