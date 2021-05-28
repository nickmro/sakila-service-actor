package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/nickmro/sakila-service-actor/sakila"
)

// ActorResolver returns an actor resolver.
func ActorResolver(service sakila.ActorService) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["actorId"].(int)
		if !ok {
			return nil, nil
		}

		return service.GetActor(p.Context, id)
	}
}

// ActorsResolver resolves actors.
func ActorsResolver(service sakila.ActorService) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		params := ActorParams(p.Args)
		return service.GetActors(p.Context, params)
	}
}

// ActorParams returns sakila actor query params.
func ActorParams(args map[string]interface{}) sakila.ActorParams {
	params := sakila.ActorParams{}

	if actorIDs, ok := args["actorIds"].([]interface{}); ok {
		params.ActorIDs = make([]int, len(actorIDs))
		for i := range actorIDs {
			params.ActorIDs[i] = actorIDs[i].(int)
		}
	}

	return params
}
