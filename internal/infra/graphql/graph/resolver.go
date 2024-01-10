package graph

import "github.com/zHenriqueGN/CentralLogger/internal/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RegisterSystemUseCase *usecase.RegisterSystemUseCase
	RegisterLogUseCase    *usecase.RegisterLogUseCase
}
