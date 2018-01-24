import React, { Component } from 'react'
import './App.css'

import { BrowserRouter as Router, Link, Route } from 'react-router-dom'

class App extends Component {
  render () {
    return (
      <Router>
        <div>
          <header>
            <h1><Link to="/">Welcome to go-piggy</Link></h1>
            <nav>
              <ul>
                <li><Link to="/new-manuscript">Create new manuscript</Link></li>
                <li><Link to="/manuscripts/cj-rules">/manuscript/cj-rules</Link></li>
              </ul>
            </nav>
          </header>

          <main>
            <Route exact path="/" component={Home}/>
            <Route path="/manuscripts/:entityId" component={ViewManuscript}/>
            <Route path="/new-manuscript" component={CreateManuscript}/>
          </main>
        </div>
      </Router>
    )
  }
}

const Home = () => (
  <div>
    <h3>Home</h3>
  </div>
)

const ViewManuscript = ({match}) => (
  <div>
    <h3>Manuscript {match.params.entityId}</h3>
  </div>
)

const CreateManuscript = ({match}) => (
  <div>
    <h3>I will create a manuscript...</h3>
  </div>
)

export default App
