package main

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/nickmro/sakila-service-actor/sakila/config"
	"github.com/nickmro/sakila-service-actor/sakila/graphql"
	"github.com/nickmro/sakila-service-actor/sakila/health"
	"github.com/nickmro/sakila-service-actor/sakila/http"
	"github.com/nickmro/sakila-service-actor/sakila/log"
	"github.com/nickmro/sakila-service-actor/sakila/mysql"
	"github.com/nickmro/sakila-service-actor/sakila/redis"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env, err := config.ReadEnv()
	if err != nil {
		panic(err)
	}

	logger, err := log.NewWriter(log.Environment(env.GetLogger()))
	if err != nil {
		panic(err)
	}

	defer logger.Flush()

	db, err := mysql.Open(env.GetMySQLURL())
	if err != nil {
		panic(err)
	}

	// nolint: errcheck
	defer db.Close()

	cache, err := redis.NewCache(&redis.ClientParams{
		Host:     env.GetRedisHost(),
		Port:     env.GetRedisPort(),
		Password: env.GetRedisPassword(),
	})
	if err != nil {
		panic(err)
	}

	// nolint: errcheck
	defer cache.Close()

	actorDB := &mysql.ActorService{
		DB:     db,
		Logger: logger,
	}

	actorCache := &redis.ActorService{
		ActorService: actorDB,
		Cache:        cache,
		Logger:       logger,
	}

	schema, err := graphql.NewSchema(actorCache)
	if err != nil {
		panic(err)
	}

	checker, err := health.NewChecker([]*health.Check{
		{
			Name:    "mysql",
			Checker: db,
		},
		{
			Name:    "redis",
			Checker: cache,
		},
	})
	if err != nil {
		panic(err)
	}

	if err := checker.Start(); err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(http.RequestLogger(logger))
	router.Mount("/graphql", graphql.NewHandler(schema))
	router.Mount("/healthz", health.NewHandler(checker))
	router.Mount("/readyz", health.NewHandler(checker))

	port := fmt.Sprintf(":%d", env.GetPort())

	fmt.Println("Listining on " + port)

	if err := http.ListenAndServe(port, router); err != nil {
		logger.Error(err)
	}
}
