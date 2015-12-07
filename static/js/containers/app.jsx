import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Chart from '../components/charts';
import * as Actions from '../actions/index';
import { Link } from 'react-router';
import { pushState } from 'redux-router';
import Loader from 'react-loader';


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
                <nav className="navbar navbar-dark">
                    <div className="container">
                        <a className="navbar-brand" href="#">Poda</a>
                        <ul className="nav navbar-nav">
                            <li className="nav-item active">
                                <Link className="nav-link" to='/'>
                                    Overview <span className="sr-only">(current)</span>
                                </Link>
                            </li>
                            <li className="nav-item">
                                <Link to='/bloat/indexes' className="nav-link">
                                    Indexes Bloat
                                </Link>
                            </li>
                            <li className="nav-item">
                                <Link to='/bloat/tables' className="nav-link">
                                    Tables Bloat
                                </Link>
                            </li>
                            <li className="nav-item">
                                <a className="nav-link" href="#">About</a>
                            </li>
                        </ul>
                    </div>
                </nav>
                <div className="container">
                    <Loader loaded={!this.props.state.app.isLoading}>
                        {this.props.children}
                    </Loader>
                </div>
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
