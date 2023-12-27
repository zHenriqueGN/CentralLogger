package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zHenriqueGN/CentralLogger/config"
	"github.com/zHenriqueGN/CentralLogger/internal/event"
	"github.com/zHenriqueGN/CentralLogger/internal/event/handler"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/grpc/pb"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/grpc/service"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/repository/postgres"
	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
	"github.com/zHenriqueGN/CentralLogger/pkg/events"
	"github.com/zHenriqueGN/UnitOfWork/uow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	envVars, err := config.LoadEnvVars()
	if err != nil {
		log.Fatalf("error on loading env vars: %w", err)
	}
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			envVars.PostgresUser,
			envVars.PostgresPassword,
			envVars.PostgresHost,
			envVars.PostgresPort,
			envVars.PostgresDB,
		),
	)
	rabbitMQConn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			envVars.RabbitMQUser,
			envVars.RabbitMQPassword,
			envVars.RabbitMQHost,
			envVars.RabbitMQPort,
		),
	)
	if err != nil {
		log.Fatalf("error on connecting to rabbitmq: %w", err)
	}
	defer rabbitMQConn.Close()
	rabbitMQChannel, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("error on opening rabbitmq channel: %w", err)
	}
	defer rabbitMQChannel.Close()
	logSaved := event.NewLogSaved()
	logSavedHandler := handler.NewLogSavedHandler(rabbitMQChannel)
	systemCreated := event.NewSystemCreated()
	systemCreatedHandler := handler.NewSystemCreatedHandler(rabbitMQChannel)
	dispatcher := events.NewDispatcher()
	dispatcher.Register(logSaved.GetName(), logSavedHandler)
	dispatcher.Register(systemCreated.GetName(), systemCreatedHandler)
	unitOfWork := uow.NewUnitOfWork(db)
	unitOfWork.Register("LogRepository", func(dbtx uow.DBTX) interface{} {
		return postgres.NewLogRepository(dbtx)
	})
	unitOfWork.Register("SystemRepository", func(dbtx uow.DBTX) interface{} {
		return postgres.NewSystemRepository(dbtx)
	})
	registerLogUseCase := usecase.NewRegisterLogUseCase(unitOfWork, logSaved, dispatcher)
	registerSystemUseCase := usecase.NewRegisterSystemUseCase(unitOfWork, systemCreated, dispatcher)
	registerLogService := service.NewLogService(registerLogUseCase)
	registerSystemService := service.NewSystemService(registerSystemUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterLogServiceServer(grpcServer, registerLogService)
	pb.RegisterSystemServiceServer(grpcServer, registerSystemService)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", envVars.GRPCServerPort))
	if err != nil {
		log.Fatalf("error on listening grpc server address: %w", err)
	}
	fmt.Printf("starting gRPC server on: localhost:%s\n", envVars.GRPCServerPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("error on serving grpc server: %w", err)
	}
}
