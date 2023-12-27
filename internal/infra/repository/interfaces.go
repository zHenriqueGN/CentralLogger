package repository

import (
	"context"

	"github.com/zHenriqueGN/CentralLogger/internal/entity"
)

type LogRepositoryInterface interface {
	Save(ctx context.Context, log *entity.Log) error
}

type SystemRepositoryInterface interface {
	Create(ctx context.Context, system *entity.System) error
}
