import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Chart from '../components/charts';
import * as Actions from '../actions/index';
import {formatPercent, formatBytes} from '../utils';


class Home extends Component {
    static propTypes = {
        actions: PropTypes.object,
        pushState: PropTypes.func.isRequired,
        state: PropTypes.object
    }
    constructor(props) {
        super(props);
    }
    render() {
        const { actions, state} = this.props;
        const indexBloat = state.app.metrics.top_index_bloat.map((v)=> v.bloat_ratio);
        const xIndexBloat = state.app.metrics.top_index_bloat.map((v)=> v.timestamp);
        const indexWaste = state.app.metrics.total_index_bloat_bytes.map((v)=> v.bloat_bytes);
        const xIndexWaste = state.app.metrics.total_index_bloat_bytes.map((v)=> v.timestamp);
        const tableBloat = state.app.metrics.top_table_bloat.map((v)=> v.bloat_ratio);
        const xTableBloat = state.app.metrics.top_table_bloat.map((v)=> v.timestamp);
        const tableWaste = state.app.metrics.total_table_bloat_bytes.map((v)=> v.bloat_bytes);
        const xTableWaste = state.app.metrics.total_table_bloat_bytes.map((v)=> v.timestamp);
        return (
            <div>
                <Chart data={indexBloat} x={xIndexBloat} yMax={100} yFormatter={formatPercent}/>
                <Chart data={indexWaste} x={xIndexWaste} yFormatter={formatBytes} />
                <Chart data={tableBloat} x={xTableBloat} yMax={100} yFormatter={formatPercent}/>
                <Chart data={tableWaste} x={xTableWaste} yFormatter={formatBytes}/>
            </div>
        );
    }
}

function mapStateToProps(state) {
    return {state};
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators(Actions, dispatch)
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(Home);
