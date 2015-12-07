import React from 'react';
import { Route, IndexRoute } from 'react-router';
import App from './containers/app';
import Home from './containers/home';
import {ListIndexBloat} from './containers/bloat';

export default (
        <Route path="/" component={App}>
            <IndexRoute component={Home} />
            <Route path="bloat/indexes" component={ListIndexBloat} />
        </Route>
)
