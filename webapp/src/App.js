import React, {Component} from 'react'
import ManuscriptAPI from './manuscript/Api.js'
import './App.css'

import {BrowserRouter as Router, Link, Route} from 'react-router-dom'
import ViewManuscript from "./manuscript/ViewManuscript";
import CreateManuscript from "./manuscript/CreateManuscript";

const api = new ManuscriptAPI(process.env.REACT_APP_PIG_URL)

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

export default App
