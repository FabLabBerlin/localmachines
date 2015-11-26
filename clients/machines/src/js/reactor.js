var Nuclear = require('nuclear-js');


export default new Nuclear.Reactor({
  debug: window.location.host !== 'easylab.io' && window.location.host !== 'lab.fablab.berlin'
});
