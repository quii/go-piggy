# go-piggy

A playground for me to mess around with CQRS and event sourcing. 

Clone into `$GOPATH/src/github.com/quii` as usual for Go projects

## dependencies

- dep `brew install dep`

`./build.sh`

## Notes

- The root package contains the main event sourcing code. (in memory source might move when > 1 implementation exists)
- Manuscript is an example package using event sourcing to contruct a projection from events. 
- cmd is an app that brings it all together as a HTTP server, you can `POST /` to create a document which returns a `Location` header for you to follow to `GET` it.

### Tentative next steps

- PATCH documents.
- Finish off versioning code, tests and then expose over HTTP
- Make a manuscript viewer. Be able to view different versions
    - Make a flashy JS version with a slider
- Metadata on events, `author`, `datetime`, `notes`    
- Create a db backed event source (maybe bolt db for fun)
- Improve naming to be more like how Young describes (https://cqrs.files.wordpress.com/2010/11/cqrs_documents.pdf)

> Command. A command is a simple object with a name of an operation and 
  the data required to perform 
  that operation. Many think of Commands as being Serializable Method Calls
  
> One important aspect of Commands is that they are always in the imperative tense; that is they are 
  telling the Application Server to do something.