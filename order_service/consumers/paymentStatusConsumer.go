package consumers

import (
	"context"
	"encoding/json"
	"log"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"example.com/rest-api/rabbitmq"
)

func PaymentStatusConsumer() {
	payment, err := rabbitmq.ConsumeEvents("exchange", "payment_status_queue", "topic", "payment.status")
	if err != nil {
		log.Fatalf("Failed to start consuming order events: %v", err)
	}

	// Fetch the payment status from the messages
	go func() {
		for msg := range payment {
			var payment_status dto.UpdatePaymentStatusDTO
			err := json.Unmarshal(msg.Body, &payment_status)
			if err != nil {
				log.Printf("Error unmarshaling order event: %v", err)
				continue
			}

			// Process the payment event
			err = handleOrderCreated(context.Background(), payment_status)
			if err != nil {
				log.Printf("Error processing order event: %v", err)
				continue
			}
			msg.Ack(false)
		}
	}()

	log.Println("Order event consumer started...")
}

func handleOrderCreated(_ context.Context, payment dto.UpdatePaymentStatusDTO) error {
	// Fetch the order from id
	order, err := models.GetOrderByID(*payment.OrderId)
	if err != nil {
		return err
	}

	// Update order status
	order.Status = *payment.PaymentStatus
	err = order.UpdateStatus()
	if err != nil {
		return err
	}
	log.Printf("Updated order status for order id: %v", order.ID)
	return nil
}
