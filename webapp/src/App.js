import React, { Component } from 'react'
import './App.css'

import { BrowserRouter as Router, Link, Redirect, Route } from 'react-router-dom'

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

class ViewManuscript extends React.Component {

  constructor(props, context) {
    super()
    this.id = props.entityId
  }

  getManuscript(id) {
    fetch(id)
      .then(res => res.json())
      .then(manuscript => this.setState({manuscript}))
  }

  render () {
    if(this.state && this.state.manuscript) {
      return (<div>
        <h1><em>Title</em> {this.state.manuscript.Title}</h1>
        <h2>Manuscript {this.entityId}</h2>
        <article>
          <h3>Abstract</h3>
          {this.state.manuscript.Abstract}
          </article>
      </div>)
    } else {
      this.getManuscript(this.entityId)
      return (<div>
      </div>)
    }
  }
}

class CreateManuscript extends React.Component {

  createNewManuscript () {
    fetch('/manuscripts', {method: 'POST'})
      .then(res => res.headers.get('location'))
      .then(location => this.setState({location}))
  }

  render () {
    if (this.state && this.state.location) {
      return <Redirect to={this.state.location}/>
    } else {
      this.createNewManuscript()
      return <p>Creating new manuscript</p>
    }
  }
}

export default App
