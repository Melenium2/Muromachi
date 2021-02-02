package graph

import "Muromachi/store"

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your apprepo, add any dependencies you require here.

type Resolver struct{
	Tables *store.TableCollection
}
