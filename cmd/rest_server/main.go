package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	dependencyconfigs "github.com/zHenriqueGN/CentralLogger/cmd/dependency_configs"
	"github.com/zHenriqueGN/CentralLogger/config"
	"github.com/zHenriqueGN/CentralLogger/internal/infra/rest/router"
	"github.com/zHenriqueGN/CentralLogger/internal/usecase"
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

	r := configureRouter(registerUseCasesOutput.RegisterSystemUseCase, registerUseCasesOutput.RegisterLogUseCase)
	http.ListenAndServe(envVars.RESTServerPort, r)
}

func configureRouter(registerSystemUseCase *usecase.RegisterSystemUseCase, registerLogUseCase *usecase.RegisterLogUseCase) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	systemRouter := router.NewSystemRouter(registerSystemUseCase)
	r.Route("/systems", func(r chi.Router) {
		r.MethodFunc(http.MethodPost, "/", systemRouter.Register)
	})
	logRouter := router.NewLogRouter(registerLogUseCase)
	r.Route("/logs", func(r chi.Router) {
		r.MethodFunc(http.MethodPost, "/", logRouter.Register)
	})
	return r
}
