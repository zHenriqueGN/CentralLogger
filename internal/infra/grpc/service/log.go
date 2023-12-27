package service

import (
	"context"
	"time"

	"github.com/zHenriqueGN/CentralLogger/internal/infra/grpc/pb"
	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
)

type LogService struct {
	pb.UnimplementedLogServiceServer
	RegisterLogUseCase *usecase.RegisterLogUseCase
}

func NewLogService(registerLogUseCase *usecase.RegisterLogUseCase) *LogService {
	return &LogService{
		RegisterLogUseCase: registerLogUseCase,
	}
}

func (l *LogService) RegisterLog(ctx context.Context, input *pb.RegisterLogRequest) (*pb.RegisterLogResponse, error) {
	inputTimeStamp, err := time.Parse(time.RFC3339, input.TimeStamp)
	if err != nil {
		return nil, err
	}
	dto := usecase.RegisterLogUseCaseInputDTO{
		SystemID:  input.SystemId,
		Level:     input.Level,
		Status:    input.Status,
		Message:   input.Message,
		TimeStamp: &inputTimeStamp,
	}
	output, err := l.RegisterLogUseCase.Execute(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterLogResponse{
		Id:        output.ID,
		SystemId:  output.SystemID,
		Level:     output.Level,
		Message:   output.Message,
		TimeStamp: output.TimeStamp.Format(time.RFC3339),
	}, nil
}
