package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"example.com/rest-api/dto"
	"example.com/rest-api/rabbitmq"
	"example.com/rest-api/utils"
)

// EmailSender allows mocking in tests
type EmailSender interface {
	Send(to, subject, body string)
}

// UserEmailFetcher allows mocking in tests
type UserEmailFetcher interface {
	GetEmail(userID string) (string, error)
}

// defaultEmailSender wraps utils.SendEmail
type defaultEmailSender struct{}

func (d *defaultEmailSender) Send(to, subject, body string) {
	utils.SendEmail(to, subject, body)
}


func OrderCreatedConsumer() {
	msgs, err := rabbitmq.ConsumeEvents("exchange", "notification_order_queue", "topic", "order.created")
	if err != nil {
		log.Fatalf("Failed to start order created consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var event dto.OrderCreatedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Error unmarshaling order.created event: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := HandleOrderCreated(event, utils.NewUserEmailFetcher(), &defaultEmailSender{}); err != nil {
				log.Printf("Error handling order.created event: %v", err)
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
		}
	}()

	log.Println("OrderCreatedConsumer started...")
}

func HandleOrderCreated(event dto.OrderCreatedEvent, fetcher UserEmailFetcher, sender EmailSender) error {
	email, err := fetcher.GetEmail(event.UserID)
	if err != nil {
		return fmt.Errorf("could not fetch user email: %w", err)
	}

	subject := "Order Confirmed!"
	body := fmt.Sprintf(
		"Hi! Your order #%s has been placed successfully on %s. Total: %.2f. Current status: %s.",
		event.ID, event.OrderDate, event.TotalAmount, event.Status,
	)

	sender.Send(email, subject, body)
	return nil
}
