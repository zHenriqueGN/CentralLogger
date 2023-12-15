package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/zHenriqueGN/CentralLogger/internal/entity"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/repository"
)

type RegisterLogUseCaseInputDTO struct {
	SystemID  string
	Level     string
	Status    string
	Message   string
	TimeStamp *time.Time
	UserID    string
}

type RegisterLogUseCaseOutputDTO struct {
	ID        uuid.UUID  `json:"id"`
	SystemID  string     `json:"system_id"`
	Level     string     `json:"level"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	TimeStamp *time.Time `json:"time_stamp"`
	UserID    string     `json:"user_id"`
}

type RegisterLogUseCase struct {
	LogRepository repository.LogRepositoryInterface
}

func NewRegisterLogUseCase(logRepository repository.LogRepositoryInterface) *RegisterLogUseCase {
	return &RegisterLogUseCase{LogRepository: logRepository}
}

func (r *RegisterLogUseCase) Execute(input RegisterLogUseCaseInputDTO) (*RegisterLogUseCaseOutputDTO, error) {
	log, err := entity.NewLog(input.SystemID, input.Level, input.Status, input.Message, input.TimeStamp, input.UserID)
	if err != nil {
		return nil, err
	}
	err = log.Validate()
	if err != nil {
		return nil, err
	}
	err = r.LogRepository.Save(log)
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
		UserID:    log.UserID,
	}
	return &output, nil
}
