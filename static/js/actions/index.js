import {REQUEST_METRICS, RECEIVE_METRICS} from './types';



/*
 * action creators
 */

function requestMetrics() {
    return {
        type: REQUEST_METRICS
    };
}

function receiveMetrics(metrics) {
    return {
        type: RECEIVE_METRICS,
        value: metrics
    };
}

export function fetchMetrics() {
    return dispatch => {
        dispatch(requestMetrics());
        return fetch("/api/stats")
            .then(response => response.json())
            .then(json => dispatch(receiveMetrics(json)));
    };
}


