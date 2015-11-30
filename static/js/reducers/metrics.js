import { RECEIVE_METRICS } from '../actions/types';

const initialState = {
    metrics: []
};

export default function metrics(state = initialState, action) {
  switch (action.type) {
    case RECEIVE_METRICS:
      return action.value;
    default:
      return state;
  }
}
