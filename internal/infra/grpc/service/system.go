package service

import (
	"context"

	"github.com/zHenriqueGN/CentralLogger/internal/infra/grpc/pb"
	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
)

type SystemService struct {
	pb.UnimplementedSystemServiceServer
	RegisterSystemUseCase *usecase.RegisterSystemUseCase
}

func NewSystemService(registerSystemUseCase *usecase.RegisterSystemUseCase) *SystemService {
	return &SystemService{
		RegisterSystemUseCase: registerSystemUseCase,
	}
}

func (s *SystemService) RegisterSystem(ctx context.Context, input *pb.RegisterSystemRequest) (*pb.RegisterSystemResponse, error) {
	dto := usecase.RegisterSystemUseCaseInputDTO{
		Name:        input.Name,
		Description: input.Description,
		Version:     input.Version,
	}
	output, err := s.RegisterSystemUseCase.Execute(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterSystemResponse{
		Id:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		Version:     output.Version,
	}, nil
}
