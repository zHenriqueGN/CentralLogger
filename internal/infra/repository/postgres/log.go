package postgres

import (
	"context"

	"github.com/zHenriqueGN/CentralLogger/internal/entity"
	"github.com/zHenriqueGN/UnitOfWork/uow"
)

type LogRepository struct {
	dbtx uow.DBTX
}

func NewLogRepository(dbtx uow.DBTX) *LogRepository {
	return &LogRepository{dbtx: dbtx}
}

func (l *LogRepository) Save(ctx context.Context, system *entity.Log) error {
	stmt, err := l.dbtx.PrepareContext(ctx, "INSERT INTO logs (id, system_id, level, status, message, time_stamp) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, system.ID, system.SystemID, system.Level, system.Status, system.Message, system.TimeStamp)
	if err != nil {
		return err
	}
	return nil
}
