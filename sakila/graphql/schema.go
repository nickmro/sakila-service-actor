package graphql

import (
	"encoding/json"

	"github.com/nickmro/sakila-service-actor/sakila"

	"github.com/graphql-go/graphql"
)

// Schema is a sakila graphQL schema.
type Schema struct {
	*graphql.Schema
}

// NewSchema returns a new graphQL schema.
func NewSchema(service sakila.ActorService) (*Schema, error) {
	actorType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Actor",
			Description: "A film actor.",
			Fields: graphql.Fields{
				"actorId": &graphql.Field{
					Type:        graphql.Int,
					Description: "The actor ID.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if actor, ok := p.Source.(*sakila.Actor); ok {
							return actor.ActorID, nil
						}

						return nil, nil
					},
				},
				"firstName": &graphql.Field{
					Type:        graphql.String,
					Description: "The actor first name.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if actor, ok := p.Source.(*sakila.Actor); ok {
							return actor.FirstName, nil
						}

						return nil, nil
					},
				},
				"lastName": &graphql.Field{
					Type:        graphql.String,
					Description: "The actor last name.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if actor, ok := p.Source.(*sakila.Actor); ok {
							return actor.LastName, nil
						}

						return nil, nil
					},
				},
				"lastUpdate": &graphql.Field{
					Type:        graphql.DateTime,
					Description: "The actor last updated at time.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if actor, ok := p.Source.(*sakila.Actor); ok {
							return actor.LastUpdate, nil
						}

						return nil, nil
					},
				},
			},
		},
	)

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Query",
					Fields: graphql.Fields{
						"actor": &graphql.Field{
							Description: "Returns an actor.",
							Type:        actorType,
							Args: graphql.FieldConfigArgument{
								"actorId": &graphql.ArgumentConfig{
									Type:        graphql.Int,
									Description: "The actor ID.",
								},
							},
							Resolve: ActorResolver(service),
						},
						"actors": &graphql.Field{
							Description: "Returns actors.",
							Type:        graphql.NewList(actorType),
							Args: graphql.FieldConfigArgument{
								"actorIds": &graphql.ArgumentConfig{
									Type:        graphql.NewList(graphql.Int),
									Description: "The actor IDs.",
								},
							},
							Resolve: ActorsResolver(service),
						},
					},
				},
			),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Schema{Schema: &schema}, nil
}

// Request takes a query to return data from the graphQL service.
func (s *Schema) Request(query string) ([]byte, error) {
	params := graphql.Params{Schema: *s.Schema, RequestString: query}

	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		return nil, r.Errors[0]
	}

	return json.Marshal(r.Data)
}
