package sakila

import (
	"context"
	"time"
)

// Actor is a sakila film actor.
type Actor struct {
	ActorID    int       `json:"actorId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	LastUpdate time.Time `json:"lastUpdate"`
}

// ActorParams are actor query parameters.
type ActorParams struct {
	ActorIDs []int
}

// ActorService defines the operations on actors.
type ActorService interface {
	GetActor(ctx context.Context, id int) (*Actor, error)
	GetActors(ctx context.Context, params ActorParams) ([]*Actor, error)
}
