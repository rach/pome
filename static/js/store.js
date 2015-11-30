import { createStore, compose } from 'redux';
import rootReducer from './reducers';
//import routes from './routes';


// const composedCreateStore = compose(
//     reduxReactRouter({ routes })
// )(createStore);

// export default function configureStore(initialState) {
//     return finalCreateStore(rootReducer, initialState);
// }

export default function configureStore(initialState) {
    const store = createStore(rootReducer, initialState);
    return store;
}
