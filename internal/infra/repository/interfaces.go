package repository

import "github.com/zHenriqueGN/CentralLogger/internal/entity"

type LogRepositoryInterface interface {
	Save(log *entity.Log) error
}

type SystemRepositoryInterface interface {
	Create(system *entity.System) error
}
