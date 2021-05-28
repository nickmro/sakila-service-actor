package graphql_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nickmro/sakila-service-actor/sakila"
	"github.com/nickmro/sakila-service-actor/sakila/graphql"
	"github.com/nickmro/sakila-service-actor/sakila/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Data struct {
	Actor  *sakila.Actor
	Actors []*sakila.Actor
}

var _ = Describe("Actor", func() {
	var actorService *mock.ActorService
	var schema *graphql.Schema

	BeforeEach(func() {
		actorService = &mock.ActorService{}

		s, err := graphql.NewSchema(actorService)
		if err != nil {
			panic(err)
		}

		schema = s

		actorService.GetActorFn = func(ctx context.Context, id int) (*sakila.Actor, error) {
			return &sakila.Actor{
				ActorID:    1,
				FirstName:  "Joe",
				LastName:   "Swank",
				LastUpdate: time.Now(),
			}, nil
		}

		actorService.GetActorsFn = func(ctx context.Context, params sakila.ActorParams) ([]*sakila.Actor, error) {
			return []*sakila.Actor{
				{
					ActorID:    1,
					FirstName:  "Joe",
					LastName:   "Swank",
					LastUpdate: time.Now(),
				},
			}, nil
		}
	})

	Describe("actor", func() {
		It("returns an actor", func() {
			query := `
				{
					actor(actorId: 1) {
						actorId
						firstName
						lastName
						lastUpdate
					}
				}
			`

			b, err := schema.Request(query)
			Expect(err).ToNot(HaveOccurred())

			data := dataFromBytes(b)
			Expect(data).ToNot(BeNil())
			Expect(data.Actor).ToNot(BeNil())
			Expect(data.Actor.ActorID).To(Equal(1))
			Expect(data.Actor.FirstName).To(Equal("Joe"))
			Expect(data.Actor.LastName).To(Equal("Swank"))
			Expect(data.Actor.LastUpdate).ToNot(BeZero())
		})
	})

	Describe("actors", func() {
		It("returns actors", func() {
			query := `
				{
					actors {
						actorId
						firstName
						lastName
						lastUpdate
					}
				}
			`

			b, err := schema.Request(query)
			Expect(err).ToNot(HaveOccurred())

			data := dataFromBytes(b)
			Expect(data).ToNot(BeNil())
			Expect(data.Actors).ToNot(BeNil())
			Expect(data.Actors).To(HaveLen(1))
			Expect(data.Actors[0].ActorID).To(Equal(1))
			Expect(data.Actors[0].FirstName).To(Equal("Joe"))
			Expect(data.Actors[0].LastName).To(Equal("Swank"))
			Expect(data.Actors[0].LastUpdate).ToNot(BeZero())
		})
	})
})

func dataFromBytes(b []byte) *Data {
	var data Data

	if err := json.Unmarshal(b, &data); err != nil {
		panic(err)
	}

	return &data
}
