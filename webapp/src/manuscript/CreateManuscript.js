import React from "react";
import ManuscriptAPI from "./Api";
import Redirect from "react-router-dom/es/Redirect";

const api = new ManuscriptAPI(process.env.REACT_APP_PIG_URL)

class CreateManuscript extends React.Component {

    createNewManuscript() {
        api.createManuscript()
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

export default CreateManuscript