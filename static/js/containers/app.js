import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Chart from '../components/charts';
import * as Actions from '../actions/index';
import { Link } from 'react-router';
import { pushState } from 'redux-router';



class App extends Component {
    propTypes = {
        actions: PropTypes.object,
        pushState: PropTypes.func.isRequired,
        state: PropTypes.object
    }
    constructor(props) {
        super(props);
    }
    render() {
        return (
            <div>
                <Link to='/'>Home</Link>
                <Link to='/bloat/indexes'>Indexes Bloat</Link>
                <Link to='/bloat/tables'>Tables Bloat</Link>
                {this.props.children}  
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

export default connect(mapStateToProps, mapDispatchToProps)(App);
