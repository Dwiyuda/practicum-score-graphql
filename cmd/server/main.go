// Command server runs the practicum score GraphQL API.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/Dwiyuda/practicum-score-graphql/graph"
	"github.com/Dwiyuda/practicum-score-graphql/internal/store"
)

func main() {
	seed := flag.Bool("seed", false, "insert synthetic demo data, then keep serving")
	flag.Parse()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	db, err := connect(ctx, url)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Migrate(context.Background()); err != nil {
		log.Fatal("migrate: ", err)
	}
	if *seed {
		if err := seedDemo(context.Background(), db); err != nil {
			log.Fatal("seed: ", err)
		}
	}

	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Store: db}}))
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.Use(extension.Introspection{})

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("Practicum Score API", "/query"))
	mux.Handle("/query", h)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })

	log.Printf("listening on :%s (playground at /)", port)
	srv := &http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

// connect retries while Postgres is still starting inside Docker Compose.
func connect(ctx context.Context, url string) (*store.Postgres, error) {
	for {
		db, err := store.NewPostgres(ctx, url)
		if err == nil {
			return db, nil
		}
		select {
		case <-ctx.Done():
			return nil, err
		case <-time.After(time.Second):
		}
	}
}
