import React, { Component, PropTypes } from "react";
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Chart from '../components/charts';
import * as Actions from '../actions/index';
import {formatPercent} from '../utils';
import filesize from 'filesize';

class ListTableBloat extends Component {
    static propTypes = {
        actions: PropTypes.object,
        state: PropTypes.object
    }
    constructor(props) {
        super(props);
    }

    render() {
      // This should be done when we reducer the state it will avoid doing it any state change
      const table_bloat = this.props.state.app.metrics.table_bloat;
      var topCharts = [];
      var bottomCharts = [];
      for (var prop in table_bloat) {
        let schema = table_bloat[prop].table_schema;
        let table = table_bloat[prop].table_name;
        let waste = table_bloat[prop].data.map((v) => v.bloat_bytes);
        let bloat = table_bloat[prop].data.map((v) => v.bloat_ratio);
        let xAxis = table_bloat[prop].data.map((v) => v.bloat_bytes);
        let chartBloat = (
          <Chart data={bloat} x={xAxis}
            yMax={100} yFormatter={formatPercent}
            title={"Bloat in"} subtitle={`Schema ${schema} Table: ${table}`}/>
        );
        let chartWaste = (
          <Chart data={waste} x={xAxis} yFormatter={filesize}
            title={"Waste in"} subtitle={`Schema: ${schema}, Table: ${table}`} />
        );
        if (prop.startsWith('pg_catalog')){
           bottomCharts.push(chartWaste);
           bottomCharts.push(<hr/>);
           bottomCharts.push(chartBloat);
           bottomCharts.push(<hr/>);
        } else {
           topCharts.push(chartWaste);
           topCharts.push(<hr/>);
           topCharts.push(chartBloat);
           topCharts.push(<hr/>);
        }
      }

        return (
            <div>
              {topCharts}
              {bottomCharts}
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

 const ltb = connect(mapStateToProps, mapDispatchToProps)(ListTableBloat);
 export {ltb as ListTableBloat}

class ListIndexBloat extends Component {
    static propTypes = {
        actions: PropTypes.object,
        state: PropTypes.object
    }
    constructor(props) {
        super(props);
    }

    render() {
      // This should be done when we reducer the state it will avoid doing it any state change
      const index_bloat = this.props.state.app.metrics.index_bloat;
      var topCharts = [];
      var bottomCharts = [];
      for (var prop in index_bloat) {
        let schema = index_bloat[prop].table_schema;
        let table = index_bloat[prop].table_name;
        let index = index_bloat[prop].index_name;
        let waste = index_bloat[prop].data.map((v) => v.bloat_bytes);
        let bloat = index_bloat[prop].data.map((v) => v.bloat_ratio);
        let xAxis = index_bloat[prop].data.map((v) => v.bloat_bytes);
        let chartBloat = (
          <Chart data={bloat} x={xAxis}
            yMax={100} yFormatter={formatPercent}
            title={"Bloat in"} subtitle={`Schema ${schema} Table: ${table}, Index: ${index} `}/>
        );
        let chartWaste = (
          <Chart data={waste} x={xAxis} yFormatter={filesize}
            title={"Waste in"} subtitle={`Schema: ${schema}, Table: ${table}, Index: ${index} `} />
        );
        if (prop.startsWith('pg_catalog')){
           bottomCharts.push(chartWaste);
           bottomCharts.push(<hr/>);
           bottomCharts.push(chartBloat);
           bottomCharts.push(<hr/>);
        } else {
           topCharts.push(chartWaste);
           topCharts.push(<hr/>);
           topCharts.push(chartBloat);
           topCharts.push(<hr/>);
        }
      }

        return (
            <div>
              {topCharts}
              {bottomCharts}
            </div>
        );
    }
}

 const lib = connect(mapStateToProps, mapDispatchToProps)(ListIndexBloat);
 export {lib as ListIndexBloat}
