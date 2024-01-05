package router

import (
	"encoding/json"
	"net/http"

	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
)

type SystemRouter struct {
	RegisterSystemUsecase *usecase.RegisterSystemUseCase
}

func NewSystemRouter(registerSystemUseCase *usecase.RegisterSystemUseCase) *SystemRouter {
	return &SystemRouter{
		RegisterSystemUsecase: registerSystemUseCase,
	}
}

func (s *SystemRouter) Register(w http.ResponseWriter, r *http.Request) {
	var registerSystemUseCaseInputDTO usecase.RegisterSystemUseCaseInputDTO
	err := json.NewDecoder(r.Body).Decode(&registerSystemUseCaseInputDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := s.RegisterSystemUsecase.Execute(r.Context(), registerSystemUseCaseInputDTO)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
