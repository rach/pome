import jQuery from 'jquery';
import {RECEIVE_METRICS} from './types';

/*
 * action creators
 */

export function fetchMetrics() {
    var result;
    jQuery.ajax({
        type: "GET",
        url: "/api/stats",
        async: false,
        dataType: "json",
        success : function(data) {
            result = data;
        }
    });
    return {type: RECEIVE_METRICS, value: result}; 
}


