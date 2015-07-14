import App from './components/App';
import Login from './components/Login';
import MachinePage from './components/MachinePage';
import React from 'react';
import Router from 'react-router';
import {DefaultRoute, Route, Routes, NotFoundRoute} from 'react-router';

/*
 * Style dependencies for webpack
 */
require('bootstrap-less');
require('../assets/less/main.less');
require('font-awesome-webpack');
require('toastr/build/toastr.min.css');

/*
 * Defined all the routes of the panel
 */
let routes = (
  <Route name="app" path="/" handler={App} >
    <Route name="machine" handler={MachinePage} />
    <Route name="login" handler={Login} />
    <DefaultRoute handler={MachinePage} />
  </Route>
);

/*
 * Render everything in the the body of index.html
 */
Router.run(routes, Router.HashLocalion, function(Handler) {
  React.render(<Handler />, document.body);
});
