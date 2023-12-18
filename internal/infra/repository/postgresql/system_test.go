package postgresql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zHenriqueGN/CentralLogger/internal/entity"
)

type SystemRepositoryTestSuite struct {
	suite.Suite
	db     *sql.DB
	system *entity.System
}

func TestSystemRepositorySuite(t *testing.T) {
	suite.Run(t, new(SystemRepositoryTestSuite))
}

func (suite *SystemRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/central_logger?sslmode=disable")
	if err != nil {
		suite.T().Fatal(err)
	}
	resetTables(suite.T(), db)
	system, err := entity.NewSystem("test system", "description", "1.0.0")
	suite.NoError(err)
	suite.db = db
	suite.system = system
}

func (suite *SystemRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *SystemRepositoryTestSuite) TestGivenAValidSystem_WhenCreatingTheSystem_ShouldCreateTheSystemCorrectly() {
	systemRepository := NewSystemRepository(suite.db)
	err := systemRepository.Create(context.Background(), suite.system)
	suite.NoError(err)
	var dbSystem entity.System
	row := suite.db.QueryRow("SELECT id, name, description, version FROM systems WHERE id = $1", suite.system.ID)
	err = row.Scan(&dbSystem.ID, &dbSystem.Name, &dbSystem.Description, &dbSystem.Version)
	suite.NoError(err)
	suite.Equal(suite.system.ID, dbSystem.ID)
	suite.Equal(suite.system.Name, dbSystem.Name)
	suite.Equal(suite.system.Description, dbSystem.Description)
	suite.Equal(suite.system.Version, dbSystem.Version)
}
