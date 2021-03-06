class ManuscriptAPI {
    constructor(baseURL) {
        this.baseURL = baseURL
    }

    getManuscript(location) {
        return fetch(`${this.baseURL}/${location}`)
            .then(res => res.json())
    }

    createManuscript() {
        return fetch(`${this.baseURL}/manuscripts`, {method: 'POST'})
            .then(checkCreated)
            .then(res => res.headers.get('location'))
    }

    updateManuscript(originalMS, newMS) {
        const eventsURL = `${this.baseURL}/manuscripts/${originalMS.EntityID}/events`;
        var facts = []

        if (originalMS.Title.trim() !== newMS.Title.trim()) {
            facts.push({Op: 'SET', Key: 'Title', Value: newMS.Title})
        }

        if (originalMS.Abstract.trim() !== newMS.Abstract.trim()) {
            facts.push({Op: 'SET', Key: 'Abstract', Value: newMS.Abstract})
        }

        return fetch(eventsURL, {
            method: 'POST',
            body: JSON.stringify(facts)
        })
    }
}

function checkCreated(res) {
    if (res.status !== 201) {
        throw new Error(`Did not get 201 from server, got ${res.status}`)
    }
    return res
}

export default ManuscriptAPI