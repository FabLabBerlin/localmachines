jest.dontMock('../helpers');
jest.dontMock('moment');

var { Month, Time } = require('../helpers');


describe('Month', function() {
  describe('getNextMonth', function() {
    it('returns Feb 2016 for Jan 2016', function() {
      var january = new Month(0, 2016);
      var february = january.getNextMonth();
      expect(january.toString()).toEqual('January');
      expect(february.toString()).toEqual('February');
    });

    it('returns Jan 2016 for Dec 2015', function() {
      var december = new Month(11, 2015);
      var january = december.getNextMonth();
      expect(december.toString()).toEqual('December');
      expect(january.toString()).toEqual('January');
    });
  });

  describe('firstDay', function() {
    it('returns 2016-07-01 for 2016-07', function() {
      var july = new Month(6, 2016);
      var yyyymmdd = july.firstDay().format('YYYY-MM-DD');
      expect(yyyymmdd).toEqual('2016-07-01');
    });
  });

  describe('lastDay', function() {
    it('returns 2016-07-31 for 2016-07', function() {
      var july = new Month(6, 2016);
      var yyyymmdd = july.lastDay().format('YYYY-MM-DD');
      expect(yyyymmdd).toEqual('2016-07-31');
    });
  });

  describe('weeks', function() {
    it('returns 4 for 2015-02', function() {
      var february = new Month(1, 2015);
      expect(february.weeks()).toEqual(4);
    });

    it('returns 5 for 2015-12', function() {
      var december = new Month(11, 2015);
      expect(december.weeks()).toEqual(5);
    });
  });
});

describe('Time', function() {
  describe('isLargerEqual', function() {
    it('returns true for 12:00 vs 12:00', function() {
      var t = new Time('12:00');
      expect(t.isLargerEqual(t)).toBe(true);
    });

    it('returns true for 12:01 vs 11:59', function() {
      var t1201 = new Time('12:01');
      var t1159 = new Time('11:59');
      expect(t1201.isLargerEqual(t1159)).toBe(true);
    });

    it('returns false for 01:11 vs 11:11', function() {
      var t0111 = new Time('01:11');
      var t1111 = new Time('11:11');
      expect(t0111.isLargerEqual(t1111)).toBe(false);
    });
  });

  describe('toInt', function() {
    it('returns 31 for 0:31', function() {
      var t = new Time('0:31');
      expect(t.toInt()).toEqual(31);
    });

    it('returns 92 for 1:32', function() {
      var t = new Time('1:32');
      expect(t.toInt()).toEqual(92);
    });
  });
});
