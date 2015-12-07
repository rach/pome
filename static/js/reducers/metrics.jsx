import { RECEIVE_METRICS, REQUEST_METRICS } from '../actions/types';

const initialState = {
    isLoading: false,
    metrics: null
};

export default function metrics(state = initialState, action) {
  switch (action.type) {
    case REQUEST_METRICS:
      return Object.assign({}, state, {
          isLoading: true
      });
    case RECEIVE_METRICS:
      return Object.assign({}, state, {
          isLoading: false,
          metrics: action.value
      });
    default:
      return state;
  }
}
