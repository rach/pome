import { createStore, compose, combineReducers, applyMiddleware } from 'redux';
import { reduxReactRouter, routerStateReducer, ReduxRouter } from 'redux-router';
import createHistory from 'history/lib/createBrowserHistory';
import rootReducer from './reducers';
import routes from './routes';
import thunk from 'redux-thunk';


const reducer = combineReducers({
    router: routerStateReducer,
    app: rootReducer
});

const composedCreateStore = compose(
    applyMiddleware(thunk),
    reduxReactRouter({ routes, createHistory })
)(createStore);

export default function configureStore(initialState) {
    return composedCreateStore(reducer, initialState);
}
