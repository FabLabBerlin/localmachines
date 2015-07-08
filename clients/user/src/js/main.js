import App from './components/App';
import Login from './components/Login';
import UserPage from './components/UserPage';
import React from 'react';
import Router from 'react-router';
import {DefaultRoute, Route, Routes, NotFoundRoute} from 'react-router';

/*
 * The style dependencies for webpack
 */
require('bootstrap-less');
require('../assets/toastr.css');
require('../assets/less/main.less');
require('font-awesome-webpack');

/*
 * Defined all the routes the router has
 */
let routes = (
  <Route name="app" path="/" handler={App} >
    <Route name="user" handler={UserPage} />
    <Route name="login" handler={Login} />
    <DefaultRoute handler={UserPage} />
  </Route>
);

/*
 * Render Everything in the body of index.html
 */
Router.run(routes, Router.HashLocation, function(Handler) {
  React.render(<Handler />, document.body);
});
