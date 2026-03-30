package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"example.com/rest-api/dto"
	"example.com/rest-api/rabbitmq"
	"example.com/rest-api/utils"
)

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

			if err := handleOrderCreated(event); err != nil {
				log.Printf("Error handling order.created event: %v", err)
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
		}
	}()

	log.Println("OrderCreatedConsumer started...")
}

func handleOrderCreated(event dto.OrderCreatedEvent) error {
	email, err := utils.GetUserEmail(event.UserID)
	if err != nil {
		return fmt.Errorf("could not fetch user email: %w", err)
	}

	subject := "Order Confirmed!"
	body := fmt.Sprintf(
		"Hi! Your order #%s has been placed successfully on %s. Total: %.2f. Current status: %s.",
		event.ID, event.OrderDate, event.TotalAmount, event.Status,
	)

	utils.SendEmail(email, subject, body)
	return nil
}
