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
    var s = String(Math.floor(d % 3600 % 60));
    if (h.length === 1) {
      h = '0' + h;
    }
    if (m.length === 1) {
      m = '0' + m;
    }
    if (s.length === 1) {
      s = '0' + s;
    }
    var str = h + ':' + m + ':' + s + ' h';

    return str;
  }
}

function formatPrice(price) {
  return (Math.round(price * 100) / 100).toFixed(2);
}

export default {
  formatDuration,
  formatPrice
};
