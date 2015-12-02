import { Route } from 'react-router';
import App from './containers/app';
import {IndexListBloat, TableListBloat} from './containers/bloat';

export default (
        <Route path="/" component={App}>
        <Route path="bloat/indexes" component={IndexListBloat} />
        <Route path="bloat/tables" component={TableListBloat} />
        </Route>
)
