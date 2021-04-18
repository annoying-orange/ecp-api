package graph

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/annoying-orange/ecp-api/graph/generated"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *sql.DB
}

// NewResolver ...
func NewResolver(db *sql.DB) (*Resolver, error) {
	return &Resolver{
		DB: db,
	}, nil
}

// Serve ...
func (r *Resolver) Serve(route string, port string) error {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: r}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				println("Websocket CheckOrigin")
				return true
			},
		},
		InitFunc: func(ctx context.Context, payload transport.InitPayload) (context.Context, error) {
			fmt.Printf("payload: %v\n", payload)

			// // // get the user from the database
			// // user := getUserByID(db, userId)
			// user := &auth.User{
			// 	ID: 2,
			// }

			// // put it in context
			// userCtx := context.WithValue(ctx, auth.UserCtxKey, user)

			// and return it so the resolvers can see it
			return nil, nil
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	router := chi.NewRouter()
	router.Handle(route, srv)
	router.Handle("/playground", playground.Handler("GraphQL playground", route))

	handler := cors.AllowAll().Handler(router)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}
