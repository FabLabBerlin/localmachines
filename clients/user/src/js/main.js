import App from './components/App';
import Login from './components/Login';
import UserPage from './components/UserPage';
import React from 'react';
import Router from 'react-router';
import {DefaultRoute, Route, Routes, NotFoundRoute} from 'react-router';

let routes = (
    <Route name="app" path="/" handler={App} >
        <Route name="user" handler={UserPage} />
        <Route name="login" handler={Login} />
        <DefaultRoute handler={UserPage} />
        <NotFoundRoute handler={UserPage} />
    </Route>
);

Router.run(routes, Router.HistoryLocation, function(Handler) {
    React.render(<Handler />, document.body);
});
