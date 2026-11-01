package graph

import "github.com/Dwiyuda/practicum-score-graphql/internal/store"

//go:generate go run github.com/99designs/gqlgen generate

// Resolver is the root resolver; every field resolver reads through Store.
type Resolver struct {
	Store store.Store
}

// ptrs converts a slice of values into the slice of pointers gqlgen expects.
func ptrs[T any](in []T) []*T {
	out := make([]*T, len(in))
	for i := range in {
		out[i] = &in[i]
	}
	return out
}
