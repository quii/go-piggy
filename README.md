# go-piggy

A playground for me to mess around with CQRS and event sourcing. 

Clone into `$GOPATH/src/github.com/quii` as usual for Go projects

`./build.sh`

## Notes

- The root package contains the main event sourcing code. (in memory source might move when > 1 implementation exists)
- Manuscript is an example package using event sourcing to contruct a projection from events. 
- cmd is an app that brings it all together as a HTTP server, you can `POST /` to create a document which returns a `Location` header for you to follow to `GET` it.

### Tentative next steps

- Create a REST-ish API for creating and updating manuscripts to see how easy (or not) the interaction is with the event sourcing
- Create a db backed event source (maybe bolt db for fun)