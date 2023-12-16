package postgresql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/zHenriqueGN/CentralLogger/internal/entity"
)

var timeFormat = "2006-01-02 15:04:05"

type LogRepositoryTestSuite struct {
	suite.Suite
	db     *sql.DB
	system *entity.System
	log    *entity.Log
}

func TestLogRepositorySuite(t *testing.T) {
	suite.Run(t, new(LogRepositoryTestSuite))
}

func (suite *LogRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/central_logger?sslmode=disable")
	if err != nil {
		suite.T().Fatal(err)
	}
	resetTables(suite.T(), db)
	system, err := entity.NewSystem("test system", "description", "1.0.0")
	suite.NoError(err)
	db.Exec("INSERT INTO systems (id, name, description, version) VALUES ($1, $2, $3, $4)", system.ID, system.Name, system.Description, system.Version)
	suite.system = system
	currentTime := time.Now()
	log, err := entity.NewLog(suite.system.ID, "DEBUG", "SUCCESS", "message", &currentTime)
	suite.NoError(err)
	suite.log = log
	suite.db = db
}

func (suite *LogRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *LogRepositoryTestSuite) TestGivenAValidLog_WhenSavingTheLog_ShouldReceiveNil() {
	logRepository := NewLogRepository(suite.db)
	err := logRepository.Save(context.Background(), suite.log)
	suite.NoError(err)
	var dbLog entity.Log
	row := suite.db.QueryRow("SELECT id, system_id, level, status, message, time_stamp FROM logs WHERE id = $1", suite.log.ID)
	err = row.Scan(&dbLog.ID, &dbLog.SystemID, &dbLog.Level, &dbLog.Status, &dbLog.Message, &dbLog.TimeStamp)
	suite.NoError(err)
	suite.Equal(suite.log.ID, dbLog.ID)
	suite.Equal(suite.log.SystemID, dbLog.SystemID)
	suite.Equal(suite.log.Level, dbLog.Level)
	suite.Equal(suite.log.Status, dbLog.Status)
	suite.Equal(suite.log.Message, dbLog.Message)
	suite.Equal(suite.log.TimeStamp.Format(timeFormat), dbLog.TimeStamp.Format(timeFormat))
}
