package handler

import (
	"context"
	"encoding/json"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zHenriqueGN/CentralLogger/pkg/events"
)

type LogSavedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewLogSavedHandler(rabbitMQChannel *amqp.Channel) *LogSavedHandler {
	return &LogSavedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *LogSavedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	jsonOutput, _ := json.Marshal(event.GetPayload())

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.PublishWithContext(
		context.Background(),
		"logs",
		"",
		false,
		false,
		message,
	)
}
