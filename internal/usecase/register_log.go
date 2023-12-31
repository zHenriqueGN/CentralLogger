package usecase

import (
	"context"
	"time"

	"github.com/zHenriqueGN/CentralLogger/internal/entity"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/repository"
	"github.com/zHenriqueGN/CentralLogger/pkg/events"
	"github.com/zHenriqueGN/UnitOfWork/uow"
)

type RegisterLogUseCaseInputDTO struct {
	SystemID  string     `json:"system_id"`
	Level     string     `json:"level"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	TimeStamp *time.Time `json:"time_stamp"`
}

type RegisterLogUseCaseOutputDTO struct {
	ID        string     `json:"id"`
	SystemID  string     `json:"system_id"`
	Level     string     `json:"level"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	TimeStamp *time.Time `json:"time_stamp"`
}

type RegisterLogUseCase struct {
	Uow        uow.UowInterface
	LogSaved   events.EventInterface
	Dispatcher events.DispatcherInterface
}

func NewRegisterLogUseCase(uow uow.UowInterface, logSaved events.EventInterface, dispatcher events.DispatcherInterface) *RegisterLogUseCase {
	return &RegisterLogUseCase{
		Uow:        uow,
		LogSaved:   logSaved,
		Dispatcher: dispatcher,
	}
}

func (r *RegisterLogUseCase) Execute(ctx context.Context, input RegisterLogUseCaseInputDTO) (*RegisterLogUseCaseOutputDTO, error) {
	log, err := entity.NewLog(input.SystemID, input.Level, input.Status, input.Message, input.TimeStamp)
	if err != nil {
		return nil, err
	}
	err = r.Uow.Do(ctx, func() error {
		logRepository, err := r.getLogRepository(ctx)
		if err != nil {
			return err
		}
		err = logRepository.Save(ctx, log)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	output := RegisterLogUseCaseOutputDTO{
		ID:        log.ID,
		SystemID:  log.SystemID,
		Level:     log.Level,
		Status:    log.Status,
		Message:   log.Message,
		TimeStamp: log.TimeStamp,
	}
	r.LogSaved.SetPayload(output)
	err = r.Dispatcher.Dispatch(r.LogSaved)
	if err != nil {
		return nil, err
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
