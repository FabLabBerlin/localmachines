/*
 * User (Page) related getters
 */
const getUser = [
  ['userStore'],
  (userStore) => {
    return userStore.get('user');
  }
];

const getBill = [
  ['userStore'],
  (userStore) => {
    return userStore.get('bill');
  }
];

const getBillMonths = [
  getUser,
  (user) => {
    var months = [];
    var created = moment(user.get('Created'));
    if (created.unix() <= 0) {
      created = moment('2015-07-01');
    }
    for (var t = created; t.isBefore(moment()); t.add(1, 'd')) {
      months.push(t.clone());
    }
    months = _.uniq(months, function(month) {
      return month.format('MMM YYYY');
    });
    return months;
  }
];

const getMemberships = [
  ['userStore'],
  (userStore) => {
    return userStore.get('memberships');
  }
];

const getMembershipsByMonth = [
  ['userStore'],
  (userStore) => {
    var byMonths = {};
    _.each(userStore.get('memberships'), function(membership) {
      var start = moment(membership.StartDate);
      var end = moment(membership.EndDate);
      for (var t = start; t.isBefore(end); t = t.add(1, 'M')) {
        var month = t.format('MMM YYYY');
        if (!byMonths[month]) {
          byMonths[month] = [];
        }
        byMonths[month].push(membership);
      }
    });
    return byMonths;
  }
];

const getMonthlyBills = [
  getBill,
  getBillMonths,
  getMembershipsByMonth,
  (bill, billMonths, membershipsByMonth) => {
    if (!bill) {
      return undefined;
    }
    var purchases = bill.Purchases;
    var purchasesByMonth = _.groupBy(purchases.Data, function(p) {
      return moment(p.TimeStart).format('MMM YYYY');
    });
    var monthlyBills = _.map(billMonths, function(m) {
      var month = m.format('MMM YYYY');
      var monthlyBill = {
        month: month,
        memberships: [],
        purchases: [],
        sums: {
          memberships: {
            priceInclVAT: 0,
            priceExclVAT: 0,
            priceVAT: 0
          },
          purchases: {
            priceInclVAT: 0,
            priceExclVAT: 0,
            priceVAT: 0,
            durations: 0
          },
          total: {}
        }
      };

      /*
       * Collect purchases and sum for the totals
       */
      _.eachRight(purchasesByMonth[month], function(purchase) {
        var timeStart = moment(purchase.TimeStart);
        var timeEnd = moment(purchase.TimeEnd);

        var duration = purchase.Quantity;
        console.log('purchase:', purchase);
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
        }

        monthlyBill.sums.durations += duration;
        var priceInclVAT = toCents(purchase.DiscountedTotal);
        var priceExclVAT = toCents(subtractVAT(purchase.DiscountedTotal));
        var priceVAT = priceInclVAT - priceExclVAT;
        monthlyBill.sums.purchases.priceInclVAT += priceInclVAT;
        monthlyBill.sums.purchases.priceExclVAT += priceExclVAT;
        monthlyBill.sums.purchases.priceVAT += priceVAT;
        monthlyBill.purchases.push({
          MachineName: purchase.Machine ? purchase.Machine.Name : 'Purchase ' + purchase.Type,
          Type: purchase.Type,
          TimeStart: timeStart,
          duration: duration,
          priceExclVAT: priceExclVAT,
          priceVAT: priceVAT,
          priceInclVAT: priceInclVAT
        });
      });

      /*
       * Collect memberships for each month
       */
      _.each(membershipsByMonth[month], function(membership) {
        var totalPrice = toCents(membership.MonthlyPrice);
        var priceExclVat = toCents(subtractVAT(membership.MonthlyPrice));
        var vat = totalPrice - priceExclVat;
        monthlyBill.memberships.push({
          startDate: moment(membership.StartDate),
          endDate: moment(membership.EndDate),
          priceExclVAT: priceExclVat,
          priceVAT: vat,
          priceInclVAT: totalPrice
        });
        monthlyBill.sums.memberships.priceInclVAT += totalPrice;
        monthlyBill.sums.memberships.priceExclVAT += priceExclVat;
        monthlyBill.sums.memberships.priceVAT += vat;
      });

      monthlyBill.sums.total = {
        priceInclVAT: monthlyBill.sums.purchases.priceInclVAT + monthlyBill.sums.memberships.priceInclVAT,
        priceExclVAT: monthlyBill.sums.purchases.priceExclVAT + monthlyBill.sums.memberships.priceExclVAT,
        priceVAT: monthlyBill.sums.purchases.priceVAT + monthlyBill.sums.memberships.priceVAT
      };

      return monthlyBill;
    });
    monthlyBills.reverse();
    return monthlyBills;
  }
];


export default {
  getUser,
  getBill,
  getBillMonths,
  getMemberships,
  getMonthlyBills
};
