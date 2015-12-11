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
      var path = this.props.state.router.location.pathname;
      return (
        <div>
          <nav className="navbar navbar-dark">
            <div className="container">
              <Link to='/' className="navbar-brand">
                Pome 
              </Link>
              <ul className="nav navbar-nav">
                <li className={path == "/bloat/indexes" ? "nav-item active" : "nav-item"}>
                  <Link to='/bloat/indexes' className="nav-link">
                    Indexes Bloat
                    {(() => {
                       if (path == '/bloat/indexes') {
                         return <span className="sr-only">(current)</span>;
                       }
                     })()}
                  </Link>
                </li>
                <li className={path == "/bloat/tables" ? "nav-item active" : "nav-item"}>
                   <Link to='/bloat/tables' className="nav-link">
                    Tables Bloat
                    {(() => {
                       if (path == '/bloat/tables') {
                         return <span className="sr-only">(current)</span>
                       }
                     })()}
                  </Link>
                </li>
              </ul>
            </div>
          </nav>
          <div className="container content">
            <Loader loaded={!this.props.state.app.isLoading}>
              {this.props.children}
            </Loader>
          </div>
          <footer className="footer">
            <div className="container">
              <div className="row">
                <div className="col-sm-6">
                  Version {this.props.state.app.isLoading? '?' : this.props.state.app.metrics.version}
                </div>
                <div className="col-sm-6 text-right">
                  © Copyright 2015, Rachid Belaid <br/>
                  Pome is licensed under the Apache License, Version 2.0
                </div>
              </div>
            </div>
          </footer>
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
