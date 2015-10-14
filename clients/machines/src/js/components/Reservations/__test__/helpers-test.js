jest.dontMock('../helpers');
jest.dontMock('moment');

var { Month } = require('../helpers');


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
