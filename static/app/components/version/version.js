'use strict';

angular.module('fabsmith.version', [
  'fabsmith.version.interpolate-filter',
  'fabsmith.version.version-directive'
])

.value('version', '0.1');
