'use strict';

describe('fabsmith.version module', function() {
  beforeEach(module('fabsmith.version'));

  describe('version service', function() {
    it('should return current version', inject(function(version) {
      expect(version).toEqual('0.1');
    }));
  });
});
