/*
 * VAT utility functions
 */
const VAT = 0.19;

function addVAT(priceExclVAT) {
  return priceExclVAT * (1 + VAT);
}

function subtractVAT(priceInclVAT) {
  return priceInclVAT / (1 + VAT);
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

export default {
  addVAT, subtractVAT, toCents, toEuro,
  formatDate
};
