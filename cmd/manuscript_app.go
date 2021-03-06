package main

import (
	"github.com/quii/go-piggy"
	"github.com/quii/go-piggy/manuscript"
	"log"
	"net/http"
)

func main() {
	eventSource := go_piggy.NewInMemoryEventSource()
	projector := manuscript.NewProjection(eventSource, nil)
	aggregate := manuscript.NewAggregate(projector, eventSource)

	server := manuscript.NewServer(projector, aggregate)

	log.Println("Listening on 8080")

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("Couldn't start server %s", err)
	}
}
