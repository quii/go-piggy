import React from "react";
import ManuscriptAPI from "./Api";

const api = new ManuscriptAPI(process.env.REACT_APP_PIG_URL)

class ViewManuscript extends React.Component {

    constructor(props, context) {
        super()
        this.location = props.match.url
        this.handleInputChange = this.handleInputChange.bind(this)
        this.editManuscript = this.editManuscript.bind(this)
    }

    getManuscript(location) {
        api.getManuscript(location)
            .then(manuscript => this.setState({
                location,
                manuscript,
                originalManuscript: Object.assign({}, manuscript)
            }))
    }

    editManuscript() {
        api.updateManuscript(this.state.originalManuscript, this.state.manuscript)
            .then(() => this.getManuscript(this.state.location))
    }

    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        const manuscript = this.state.manuscript;
        manuscript[name] = value

        this.setState({manuscript});
    }

    render() {
        if (this.state && this.state.manuscript) {
            return <div>
                <h2>{this.location} : Version {this.state.manuscript.Version}</h2>
                <label>Title</label><input onChange={this.handleInputChange} name='Title' type='text'
                                           value={this.state.manuscript.Title}/>
                <article>
                    <h3>Abstract</h3>
                    <textarea onChange={this.handleInputChange} rows='20' cols='60'
                              name='Abstract'>{this.state.manuscript.Abstract}</textarea>
                </article>
                <button onClick={this.editManuscript}>Save changes into new version</button>
            </div>
        } else {
            this.getManuscript(this.location)
            return (<div>
            </div>)
        }
    }
}

export default ViewManuscript