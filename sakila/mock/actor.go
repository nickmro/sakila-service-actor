package mock

import (
	"context"

	"github.com/nickmro/sakila-service-actor/sakila"
)

// ActorService is a mock actor service.
type ActorService struct {
	GetActorsFn func(ctx context.Context, params sakila.ActorParams) ([]*sakila.Actor, error)
	GetActorFn  func(ctx context.Context, id int) (*sakila.Actor, error)
}

// GetActors is a mock function that returns actors.
func (service *ActorService) GetActors(ctx context.Context, params sakila.ActorParams) ([]*sakila.Actor, error) {
	if fn := service.GetActorsFn; fn != nil {
		return fn(ctx, params)
	}

	return []*sakila.Actor{}, nil
}

// GetActor is a mock function that return an actor.
func (service *ActorService) GetActor(ctx context.Context, id int) (*sakila.Actor, error) {
	if fn := service.GetActorFn; fn != nil {
		return fn(ctx, id)
	}

	return &sakila.Actor{}, nil
}
