import '../scss/app';
import React from 'react';
import Root from './containers/root';
import {fetchMetrics} from './actions';
import configureStore from './store';

const store = configureStore({metrics:[]});
store.dispatch(fetchMetrics());
window.setInterval(() => {store.dispatch(fetchMetrics())}, 5000);
React.render(<Root store={store} />, document.querySelector('#chart'));
