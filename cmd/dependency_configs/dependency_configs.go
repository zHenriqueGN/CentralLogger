package dependencyconfigs

import (
	"database/sql"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zHenriqueGN/CentralLogger/internal/event"
	"github.com/zHenriqueGN/CentralLogger/internal/event/handler"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/repository/postgres"
	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
	"github.com/zHenriqueGN/CentralLogger/pkg/events"
	"github.com/zHenriqueGN/UnitOfWork/uow"
)

type RegisterUseCasesInput struct {
	UnitOfWork    uow.UowInterface
	Dispatcher    events.DispatcherInterface
	SystemCreated *event.SystemCreated
	LogSaved      *event.LogSaved
}

type RegisterUseCasesOutput struct {
	RegisterSystemUseCase *usecase.RegisterSystemUseCase
	RegisterLogUseCase    *usecase.RegisterLogUseCase
}

type RegisterEventsOutput struct {
	Dispatcher    events.DispatcherInterface
	SystemCreated *event.SystemCreated
	LogSaved      *event.LogSaved
}

func RegisterUseCases(input *RegisterUseCasesInput) *RegisterUseCasesOutput {
	RegisterSystemUseCase := usecase.NewRegisterSystemUseCase(input.UnitOfWork, input.SystemCreated, input.Dispatcher)
	RegisterLogUseCase := usecase.NewRegisterLogUseCase(input.UnitOfWork, input.LogSaved, input.Dispatcher)
	return &RegisterUseCasesOutput{
		RegisterSystemUseCase: RegisterSystemUseCase,
		RegisterLogUseCase:    RegisterLogUseCase,
	}
}

func CreateUnitOfWork(db *sql.DB) uow.UowInterface {
	unitOfWork := uow.NewUnitOfWork(db)
	unitOfWork.Register("LogRepository", func(dbtx uow.DBTX) interface{} {
		return postgres.NewLogRepository(dbtx)
	})
	unitOfWork.Register("SystemRepository", func(dbtx uow.DBTX) interface{} {
		return postgres.NewSystemRepository(dbtx)
	})
	return unitOfWork
}

func RegisterEvents(rabbitMQChannel *amqp.Channel) *RegisterEventsOutput {
	logSaved := event.NewLogSaved()
	systemCreated := event.NewSystemCreated()

	logSavedHandler := handler.NewLogSavedHandler(rabbitMQChannel)
	systemCreatedHandler := handler.NewSystemCreatedHandler(rabbitMQChannel)

	dispatcher := events.NewDispatcher()
	dispatcher.Register(logSaved.GetName(), logSavedHandler)
	dispatcher.Register(systemCreated.GetName(), systemCreatedHandler)
	return &RegisterEventsOutput{
		Dispatcher:    dispatcher,
		LogSaved:      logSaved,
		SystemCreated: systemCreated,
	}
}

func ConnectToPostgres(user, password, host, port, database string) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			database,
		),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetRabbitMQConn(user, password, host, port string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			user,
			password,
			host,
			port,
		),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetRabbitMQChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}
