//go:build js && wasm

// Command wasm runs the GraphQL API inside the browser for the public demo.
// It serves the same schema and resolvers as cmd/server, backed by the in-memory store.
package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall/js"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/Dwiyuda/practicum-score-graphql/graph"
	"github.com/Dwiyuda/practicum-score-graphql/internal/store"
)

func main() {
	db := store.NewMemory()
	if err := store.Seed(context.Background(), db); err != nil {
		panic(err)
	}

	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Store: db}}))
	h.AddTransport(transport.POST{})
	h.Use(extension.Introspection{})

	// runGraphQL takes a JSON request body and returns a promise for the JSON response.
	js.Global().Set("runGraphQL", js.FuncOf(func(_ js.Value, args []js.Value) any {
		body := args[0].String()
		return js.Global().Get("Promise").New(js.FuncOf(func(_ js.Value, p []js.Value) any {
			// The handler starts goroutines and waits for them, which would deadlock on the JS event loop.
			go func() {
				req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, req)
				p[0].Invoke(rec.Body.String())
			}()
			return nil
		}))
	}))
	js.Global().Call("onGraphQLReady")
	select {}
}
