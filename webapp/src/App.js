import React, {Component} from 'react'
import './App.css'

import {BrowserRouter as Router, Link, Redirect, Route} from 'react-router-dom'

const baseURL = process.env.REACT_APP_PIG_URL;

class App extends Component {
    render() {
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
                        <p>API URL: {baseURL}</p>
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
        console.log(props)
        console.log(`getting ${props.match.params.entityId}`)
        this.location = props.match.url
    }

    getManuscript(location) {
        fetch(`${baseURL}/${location}`)
            .then(res => res.json())
            .then(manuscript => this.setState({manuscript}))
    }

    render() {
        if (this.state && this.state.manuscript) {
            return (<div>
                <h1><em>Title</em> {this.state.manuscript.Title}</h1>
                <h2>{this.location}</h2>
                <article>
                    <h3>Abstract</h3>
                    {this.state.manuscript.Abstract}
                </article>
            </div>)
        } else {
            this.getManuscript(this.location)
            return (<div>
            </div>)
        }
    }
}

class CreateManuscript extends React.Component {

    createNewManuscript() {
        fetch(`${baseURL}/manuscripts`, {method: 'POST'})
            .then(checkCreated)
            .then(res => res.headers.get('location'))
            .then(location => this.setState({location}))
    }

    render() {
        if (this.state && this.state.location) {
            return <Redirect to={this.state.location}/>
        } else {
            this.createNewManuscript()
            return <p>Creating new manuscript</p>
        }
    }
}

function checkCreated(res) {
    if(res.status!==201){
        throw new Error(`Did not get 201 from server, got ${res.status}`)
    }
    return res
}

export default App
