package handler

import (
	"context"
	"encoding/json"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zHenriqueGN/CentralLogger/pkg/events"
)

type SystemCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewSystemCreatedHandler(rabbitMQChannel *amqp.Channel) *SystemCreatedHandler {
	return &SystemCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *SystemCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	jsonOutput, _ := json.Marshal(event.GetPayload())

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.PublishWithContext(
		context.Background(),
		"amq.direct",
		"",
		false,
		false,
		message,
	)
}
