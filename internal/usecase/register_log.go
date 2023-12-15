package usecase

import (
	"context"
	"time"

	"github.com/zHenriqueGN/CentralLogger/internal/entity"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/repository"
	"github.com/zHenriqueGN/UnitOfWork/uow"
)

type RegisterLogUseCaseInputDTO struct {
	SystemID  string     `json:"system_id"`
	Level     string     `json:"level"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	TimeStamp *time.Time `json:"time_stamp"`
	UserID    string     `json:"user_id"`
}

type RegisterLogUseCaseOutputDTO struct {
	ID        string     `json:"id"`
	SystemID  string     `json:"system_id"`
	Level     string     `json:"level"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	TimeStamp *time.Time `json:"time_stamp"`
	UserID    string     `json:"user_id"`
}

type RegisterLogUseCase struct {
	Uow uow.UowInterface
}

func NewRegisterLogUseCase(uow uow.UowInterface) *RegisterLogUseCase {
	return &RegisterLogUseCase{Uow: uow}
}

func (r *RegisterLogUseCase) Execute(ctx context.Context, input RegisterLogUseCaseInputDTO) (*RegisterLogUseCaseOutputDTO, error) {
	log, err := entity.NewLog(input.SystemID, input.Level, input.Status, input.Message, input.TimeStamp, input.UserID)
	if err != nil {
		return nil, err
	}
	logRepository, err := r.getLogRepository(ctx)
	if err != nil {
		return nil, err
	}
	err = logRepository.Save(log)
	if err != nil {
		return nil, err
	}
	output := RegisterLogUseCaseOutputDTO{
		ID:        log.ID.String(),
		SystemID:  log.SystemID,
		Level:     log.Level,
		Status:    log.Status,
		Message:   log.Message,
		TimeStamp: log.TimeStamp,
		UserID:    log.UserID,
	}
	return &output, nil
}

func (r *RegisterLogUseCase) getLogRepository(ctx context.Context) (repository.LogRepositoryInterface, error) {
	repo, err := r.Uow.GetRepository(ctx, "LogRepository")
	if err != nil {
		return nil, err
	}
	return repo.(repository.LogRepositoryInterface), nil
}
