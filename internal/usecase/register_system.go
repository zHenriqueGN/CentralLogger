package usecase

import (
	"context"

	"github.com/zHenriqueGN/CentralLogger/internal/entity"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/repository"
	"github.com/zHenriqueGN/UnitOfWork/uow"
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
	Uow uow.UowInterface
}

func NewRegisterSystemUseCase(uow uow.UowInterface) *RegisterSystemUseCase {
	return &RegisterSystemUseCase{Uow: uow}
}

func (r *RegisterSystemUseCase) Execute(ctx context.Context, input RegisterSystemUseCaseInputDTO) (*RegisterSystemUseCaseOutputDTO, error) {
	system, err := entity.NewSystem(input.Name, input.Description, input.Version)
	if err != nil {
		return nil, err
	}
	systemRepository, err := r.getSystemRepository(ctx)
	if err != nil {
		return nil, err
	}
	err = systemRepository.Create(system)
	if err != nil {
		return nil, err
	}
	output := RegisterSystemUseCaseOutputDTO{
		ID:          system.ID,
		Name:        system.Name,
		Description: system.Description,
		Version:     system.Version,
	}
	return &output, nil
}

func (r *RegisterSystemUseCase) getSystemRepository(ctx context.Context) (repository.SystemRepositoryInterface, error) {
	repo, err := r.Uow.GetRepository(ctx, "SystemRepository")
	if err != nil {
		return nil, err
	}
	return repo.(repository.SystemRepositoryInterface), nil
}
