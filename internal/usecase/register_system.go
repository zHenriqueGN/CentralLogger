package usecase

import (
	"github.com/zHenriqueGN/CentralLogger/internal/entity"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/repository"
)

type RegisterSystemUseCaseInputDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type RegisterSystemUseCaseOutputDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type RegisterSystemUseCase struct {
	SystemRepository repository.SystemRepositoryInterface
}

func NewRegisterSystemUseCase(systemRepository repository.SystemRepositoryInterface) *RegisterSystemUseCase {
	return &RegisterSystemUseCase{SystemRepository: systemRepository}
}

func (r *RegisterSystemUseCase) Execute(input RegisterSystemUseCaseInputDTO) (*RegisterSystemUseCaseOutputDTO, error) {
	system, err := entity.NewSystem(input.Name, input.Description, input.Version)
	if err != nil {
		return nil, err
	}
	err = r.SystemRepository.Create(system)
	if err != nil {
		return nil, err
	}
	output := RegisterSystemUseCaseOutputDTO{
		ID:          system.ID.String(),
		Name:        system.Name,
		Description: system.Description,
		Version:     system.Version,
	}
	return &output, nil
}
