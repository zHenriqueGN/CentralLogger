package postgres

import (
	"context"

	"github.com/zHenriqueGN/CentralLogger/internal/entity"
	"github.com/zHenriqueGN/UnitOfWork/uow"
)

type SystemRepository struct {
	dbtx uow.DBTX
}

func NewSystemRepository(dbtx uow.DBTX) *SystemRepository {
	return &SystemRepository{dbtx: dbtx}
}

func (s *SystemRepository) Create(ctx context.Context, system *entity.System) error {
	stmt, err := s.dbtx.PrepareContext(ctx, "INSERT INTO systems (id, name, description, version) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, system.ID, system.Name, system.Description, system.Version)
	if err != nil {
		return err
	}
	return nil
}
