package redis

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/nickmro/sakila-service-actor/sakila"
)

// ActorService is a cached wrapper for an actor service.
type ActorService struct {
	sakila.ActorService
	Cache  *Cache
	TTL    time.Duration
	Logger sakila.Logger
}

// GetActors returns actors from the cache.
func (service *ActorService) GetActors(ctx context.Context, params sakila.ActorParams) ([]*sakila.Actor, error) {
	var actors []*sakila.Actor

	item := &cache.Item{
		Ctx:   ctx,
		Key:   cacheKey("actors::", params),
		Value: &actors,
		Do: func(i *cache.Item) (interface{}, error) {
			return service.ActorService.GetActors(ctx, params)
		},
		TTL: service.TTL,
	}

	err := service.Cache.Once(item)
	if err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

	return actors, nil
}

// GetActor returns an actor from the cache.
func (service *ActorService) GetActor(ctx context.Context, id int) (*sakila.Actor, error) {
	var actor sakila.Actor

	item := &cache.Item{
		Ctx:   ctx,
		Key:   cacheKey("actor::", sakila.ActorParams{ActorIDs: []int{id}}),
		Value: &actor,
		Do: func(i *cache.Item) (interface{}, error) {
			return service.ActorService.GetActor(ctx, id)
		},
		TTL: service.TTL,
	}

	err := service.Cache.Once(item)
	if err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

	return &actor, nil
}

func (service *ActorService) logError(err error) {
	if logger := service.Logger; logger != nil {
		logger.Error(err)
	}
}

func cacheKey(prefix string, params sakila.ActorParams) string {
	key := strings.Builder{}
	key.WriteString(prefix)

	if len(params.ActorIDs) > 0 {
		ids := []string{}
		for _, id := range params.ActorIDs {
			ids = append(ids, strconv.Itoa(id))
		}

		key.WriteString("id:")
		key.WriteString(strings.Join(ids, ","))
	}

	return hashedKey(key.String())
}
