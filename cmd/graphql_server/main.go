package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	dependencyconfigs "github.com/zHenriqueGN/CentralLogger/cmd/dependency_configs"
	"github.com/zHenriqueGN/CentralLogger/config"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/graphql/graph"
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

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					RegisterSystemUseCase: registerUseCasesOutput.RegisterSystemUseCase,
					RegisterLogUseCase:    registerUseCasesOutput.RegisterLogUseCase,
				},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", envVars.GraphQLServerPort)
	log.Fatal(http.ListenAndServe(":"+envVars.GraphQLServerPort, nil))
}
