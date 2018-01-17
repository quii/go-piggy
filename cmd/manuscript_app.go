package main

import (
	"github.com/quii/go-piggy"
	"github.com/quii/go-piggy/manuscript"
	"log"
	"net/http"
)

func main() {
	eventSource := go_piggy.NewInMemoryEventSource()

	projector := manuscript.NewProjection(eventSource)
	emitter := manuscript.NewEmitter(eventSource)

	server := manuscript.Server{
		Repo: struct {
			*manuscript.Emitter
			*manuscript.Projection
		}{emitter, projector},
		EntityIdGenerator: go_piggy.RandomID,
	}

	log.Println("Listening on 8080")

	if err := http.ListenAndServe(":8080", &server); err != nil {
		log.Fatalf("Couldn't start server %s", err)
	}
}
