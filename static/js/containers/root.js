import React, { Component, PropTypes } from 'react';
import { Provider } from 'react-redux';
//import { ReduxRouter } from 'redux-router';
import App from './app';

export default class Root extends Component {
    constructor(props) {
        super(props);
    }
    propTypes = {
        store: PropTypes.object.isRequired
    }
    render() {
        const { store } = this.props;
        return (
                <Provider store={store}>
                    <App />
                </Provider>
        );
    }
}

