import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { BarChart, BytesBarChart, PercentBarChart } from '../components/charts';
import * as Actions from '../actions/index';
import filesize from 'filesize';


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
        // This should be done when we reduce the state it will avoid doing it any state
        // change
        const {state} = this.props;
      const indexBloat = state.app.metrics.top_index_bloat.map((v)=> {
        return {
        x: v.timestamp,
        y: v.bloat_ratio,
        context: "Schema: " + v.table_schema + ", Table: " + v.table_name + ", Index: " + v.index_name
        };
      });
      const indexWaste = state.app.metrics.total_index_bloat_bytes.map((v)=> {
        return {
          x: v.timestamp,
          y: v.bloat_bytes
        };
      });
      const tableBloat = state.app.metrics.top_table_bloat.map((v)=> {
        return {
          x: v.timestamp,
          y: v.bloat_ratio,
          context: "Schema: " + v.table_schema + ", Table: " + v.table_name
        };
      });
      const tableWaste = state.app.metrics.total_table_bloat_bytes.map((v)=> {
        return {
          x: v.timestamp,
          y: v.bloat_bytes
        };
      });

      const databaseSize = state.app.metrics.database_size.map((v)=> {
        return {
          x: v.timestamp,
          y: v.total_size
        };
      });
      const tablesSize = state.app.metrics.database_size.map((v)=> {
        return {
          x: v.timestamp,
          y: v.table_size
        };
      });
      const indexesSize = state.app.metrics.database_size.map((v)=> {
        return {
          x: v.timestamp,
          y: v.index_size
        };
      });
      const indexesRatio = state.app.metrics.database_size.map((v)=> {
        return {
          x: v.timestamp,
          y: v.index_ratio
        };
      });
      const numOfConnection = state.app.metrics.number_of_connection.map((v)=> {
        return {
          x: v.timestamp,
          y: v.count
        };
      });
      const transactionPerSec = state.app.metrics.transaction_per_sec.map((v)=> {
        return {
          x: v.timestamp,
          y: v.tps
        };
      });
        return (
            <div>
                <BarChart data={numOfConnection} title={"Number of connections"}/>
                <hr/>
                <BarChart data={transactionPerSec} title={"Transactions per second"}/>
                <hr/>
                <PercentBarChart data={indexBloat} title={"Most bloated index"}/>
                <hr/>
                <BytesBarChart data={indexWaste} title={"Total wasted bytes for indexes"} />
                <hr/>
                <PercentBarChart data={tableBloat} title={"Most bloated table"}/>
                <hr/>
                <BytesBarChart data={tableWaste} title={"Total wasted bytes for tables"}/>
                <hr/>
                <BytesBarChart data={databaseSize} title={"Database Size"}/>
                <hr/>
                <BytesBarChart data={tablesSize} title={"Table Size"}/>
                <hr/>
                <BytesBarChart data={indexesSize} title={"Indexes Size"}/>
                <hr/>
                <PercentBarChart data={indexesRatio} title={"Indexes Size Ratio"}/>
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
