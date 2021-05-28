package mysql

import (
	"context"

	"github.com/nickmro/mrqb"
	"github.com/nickmro/sakila-service-actor/sakila"
)

// ActorService is an actor service backed by a MySQL DB.
type ActorService struct {
	*DB
	Logger sakila.Logger
}

// GetActor returns an actor.
func (service *ActorService) GetActor(ctx context.Context, id int) (*sakila.Actor, error) {
	var actor sakila.Actor

	query, args := actorsQuery(sakila.ActorParams{ActorIDs: []int{id}})

	if err := service.DB.QueryRow(query, args...).Scan(
		&actor.ActorID,
		&actor.FirstName,
		&actor.LastName,
		&actor.LastUpdate,
	); err != nil {
		service.logError(err)
		return nil, err
	}

	return &actor, nil
}

// GetActors returns actors.
func (service *ActorService) GetActors(ctx context.Context, params sakila.ActorParams) ([]*sakila.Actor, error) {
	actors := []*sakila.Actor{}

	query, args := actorsQuery(params)

	rows, err := service.DB.Query(query, args...)
	if err != nil {
		service.logError(err)
		return nil, err
	} else if rows.Err() != nil {
		service.logError(err)
		return nil, rows.Err()
	}

	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var actor sakila.Actor

		err := rows.Scan(
			&actor.ActorID,
			&actor.FirstName,
			&actor.LastName,
			&actor.LastUpdate,
		)
		if err != nil {
			service.logError(err)
			return nil, err
		}

		actors = append(actors, &actor)
	}

	return actors, nil
}

func actorsQuery(params sakila.ActorParams) (string, []interface{}) {
	query := mrqb.Select(
		"actor.actor_id",
		"actor.first_name",
		"actor.last_name",
		"actor.last_update",
	).
		From("actor")

	ids := params.ActorIDs
	if len(ids) == 1 {
		query.Where("actor.actor_id = %v", ids[0])
	} else if len(ids) > 0 {
		args := make([]interface{}, len(ids))
		for i := range ids {
			args[i] = ids[i]
		}

		query.Where("actor.actor_id IN (%v)", args...)
	}

	return query.Build()
}

func (service *ActorService) logError(err error) {
	if logger := service.Logger; logger != nil {
		logger.Error(err)
	}
}
