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
}

func (suite *LogRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/central_logger?sslmode=disable")
	suite.NoError(err)
	db.Exec("DROP TABLE IF EXISTS logs;")
	db.Exec("DROP TABLE IF EXISTS systems;")
	db.Exec(`
		CREATE TABLE systems (
    		id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY,
        	name VARCHAR(50) NOT NULL,
        	description VARCHAR(200) NOT NULL,
        	version VARCHAR(15) NOT NULL
		);
	`)
	db.Exec(`
		CREATE TABLE logs (
        	id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY,
        	system_id VARCHAR(36) NOT NULL REFERENCES systems (id),
        	level VARCHAR(30) NOT NULL,
        	status VARCHAR(30) NOT NULL,
        	message VARCHAR NOT NULL,
        	time_stamp TIMESTAMP NOT NULL
		);
	`)
	system, err := entity.NewSystem("test system", "description", "1.0.0")
	suite.NoError(err)
	db.Exec("INSERT INTO systems (id, name, description, version) VALUES ($1, $2, $3, $4)", system.ID, system.Name, system.Description, system.Version)
	suite.system = system
	suite.db = db
}

func (suite *LogRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func TestLogRepositorySuite(t *testing.T) {
	suite.Run(t, new(LogRepositoryTestSuite))
}

func (suite *LogRepositoryTestSuite) TestGivenAValidLog_WhenSavingTheLog_ShouldReceiveNil() {
	currentTime := time.Now()
	log, err := entity.NewLog(suite.system.ID, "DEBUG", "SUCCESS", "message", &currentTime)
	suite.NoError(err)
	logRepository := NewLogRepository(suite.db)
	err = logRepository.Save(context.Background(), log)
	suite.NoError(err)
	var dbLog entity.Log
	row := suite.db.QueryRow("SELECT id, system_id, level, status, message, time_stamp FROM logs WHERE id = $1", log.ID)
	err = row.Scan(&dbLog.ID, &dbLog.SystemID, &dbLog.Level, &dbLog.Status, &dbLog.Message, &dbLog.TimeStamp)
	suite.NoError(err)
	suite.Equal(log.ID, dbLog.ID)
	suite.Equal(log.SystemID, dbLog.SystemID)
	suite.Equal(log.Level, dbLog.Level)
	suite.Equal(log.Status, dbLog.Status)
	suite.Equal(log.Message, dbLog.Message)
	suite.Equal(log.TimeStamp.Format(timeFormat), dbLog.TimeStamp.Format(timeFormat))
}
