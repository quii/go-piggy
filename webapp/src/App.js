import React, { Component } from 'react';
import './App.css';

import {
  BrowserRouter as Router,
  Route,
  Link
} from 'react-router-dom'

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
        <h1 className="App-title">Welcome to go-piggy</h1>
        </header>
        <main>
          <Router>
            <div>
              <Route exact path="/" component={Home}/>
              <Route path="/manuscripts/:entityId" component={Manuscript}/>
            </div>
          </Router>
        </main>
      </div>
    );
  }
}

const Home = () => (
  <div>
    <h3>Home</h3>
    <p>Try <Link to="/manuscripts/cj-rules">/manuscript/cj-rules</Link></p>
  </div>
)

const Manuscript = ({ match }) => (
  <div>
    <h3>Manuscript {match.params.entityId}</h3>
  </div>
)

export default App;
