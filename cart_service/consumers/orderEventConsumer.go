package consumers

import (
	"context"
	"encoding/json"
	"log"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"example.com/rest-api/rabbitmq"
)

func OrderEventConsumer() {
	order, err := rabbitmq.ConsumeEvents("exchange", "order", "topic", "order.created")
	if err != nil {
		log.Fatalf("Failed to start consuming order events: %v", err)
	}

	// Fetch the order details from the messages
	go func() {
		for msg := range order {
			var orderEvent dto.OrderEvent
			err := json.Unmarshal(msg.Body, &orderEvent)
			if err != nil {
				log.Printf("Error unmarshaling order event: %v", err)
				continue
			}

			// Process the order event
			err = handleOrderCreated(context.Background(), orderEvent)
			if err != nil {
				log.Printf("Error processing order event: %v", err)
				continue
			}
			msg.Ack(false)
		}
	}()

	log.Println("Order event consumer started...")
}

func handleOrderCreated(_ context.Context, orderEvent dto.OrderEvent) error {
	// Fetch the user's active cart
	userCart, err := models.GetActiveCartByUserId(orderEvent.UserId)
	if err != nil {
		return err
	}

	// Mark the cart as inactive
	err = models.DeactivateCart(userCart.ID)
	if err != nil {
		return err
	}

	// Here you can add additional logic, such as sending a confirmation email, etc.
	log.Printf("Order processed for user %s, cart %s marked as inactive.", orderEvent.UserId, userCart.ID)
	return nil
}
