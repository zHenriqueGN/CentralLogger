package router

import (
	"encoding/json"
	"net/http"

	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
)

type LogRouter struct {
	RegisterLogUsecase *usecase.RegisterLogUseCase
}

func NewLogRouter(registerLogUseCase *usecase.RegisterLogUseCase) *LogRouter {
	return &LogRouter{
		RegisterLogUsecase: registerLogUseCase,
	}
}

func (l *LogRouter) Register(w http.ResponseWriter, r *http.Request) {
	var registerLogUseCaseInputDTO usecase.RegisterLogUseCaseInputDTO
	err := json.NewDecoder(r.Body).Decode(&registerLogUseCaseInputDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := l.RegisterLogUsecase.Execute(r.Context(), registerLogUseCaseInputDTO)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
