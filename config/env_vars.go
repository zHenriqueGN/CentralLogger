package config

import "github.com/spf13/viper"

type EnvVars struct {
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	GRPCServerPort   string `mapstructure:"GRPC_SERVER_PORT"`
	RabbitMQHost     string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQUser     string `mapstructure:"RABBITMQ_USER"`
	RabbitMQPassword string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitMQPort     string `mapstructure:"RABBITMQ_PORT"`
}

func LoadEnvVars() (*EnvVars, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var envVars EnvVars
	err = viper.Unmarshal(&envVars)
	if err != nil {
		return nil, err
	}
	return &envVars, nil
}
