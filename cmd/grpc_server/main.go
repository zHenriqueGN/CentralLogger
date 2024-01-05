package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	dependencyconfigs "github.com/zHenriqueGN/CentralLogger/cmd/dependency_configs"
	"github.com/zHenriqueGN/CentralLogger/config"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/grpc/pb"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/grpc/service"
	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	envVars, err := config.LoadEnvVars()
	if err != nil {
		log.Fatalf("error on loading env vars: %v", err)
	}

	db, err := dependencyconfigs.ConnectToPostgres(envVars.PostgresUser, envVars.PostgresPassword, envVars.PostgresHost, envVars.PostgresPort, envVars.PostgresDB)
	if err != nil {
		log.Fatalf("error when connecting to postgres: %v", err)
	}
	defer db.Close()

	rabbitMQConn, err := dependencyconfigs.GetRabbitMQConn(envVars.RabbitMQUser, envVars.RabbitMQPassword, envVars.RabbitMQHost, envVars.RabbitMQPort)
	if err != nil {
		log.Fatalf("error when connecting to RabbitMQ: %v", err)
	}
	defer rabbitMQConn.Close()

	rabbitMQChannel, err := dependencyconfigs.GetRabbitMQChannel(rabbitMQConn)
	if err != nil {
		log.Fatalf("error when getting RabbitMQ channel: %v", err)
	}
	defer rabbitMQChannel.Close()

	registerEventsOutput := dependencyconfigs.RegisterEvents(rabbitMQChannel)
	unitOfWork := dependencyconfigs.CreateUnitOfWork(db)
	registerUseCasesOutput := dependencyconfigs.RegisterUseCases(&dependencyconfigs.RegisterUseCasesInput{
		UnitOfWork:    unitOfWork,
		Dispatcher:    registerEventsOutput.Dispatcher,
		SystemCreated: registerEventsOutput.SystemCreated,
		LogSaved:      registerEventsOutput.LogSaved,
	})

	grpcServer := configuregRPCServices(registerUseCasesOutput.RegisterSystemUseCase, registerUseCasesOutput.RegisterLogUseCase)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", envVars.GRPCServerPort))
	if err != nil {
		log.Fatalf("error on listening grpc server address: %v", err)
	}
	fmt.Printf("starting gRPC server on: localhost:%s\n", envVars.GRPCServerPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("error on serving grpc server: %v", err)
	}
}

func configuregRPCServices(registerSystemUseCase *usecase.RegisterSystemUseCase, registerLogUseCase *usecase.RegisterLogUseCase) *grpc.Server {
	registerSystemService := service.NewSystemService(registerSystemUseCase)
	registerLogService := service.NewLogService(registerLogUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterSystemServiceServer(grpcServer, registerSystemService)
	pb.RegisterLogServiceServer(grpcServer, registerLogService)
	reflection.Register(grpcServer)
	return grpcServer
}
