package consumers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"example.com/rest-api/rabbitmq"
)

func OrderEventConsumer() {
	order, err := rabbitmq.ConsumeEvents("exchange", "order_payment_queue", "topic", "order.created")
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
			err = handleOrderEvent(context.Background(), orderEvent)
			if err != nil {
				log.Printf("Error processing order event: %v", err)
				continue
			}
			msg.Ack(false)
		}
	}()

	log.Println("Order event consumer started...")
}

func handleOrderEvent(_ context.Context, orderEvent dto.OrderEvent) error {

	var payment models.Payment
	payment.OrderId = orderEvent.OrderId
	payment.UserId = orderEvent.UserId
	payment.CartId = orderEvent.CartId

	delay := time.Duration(2+rand.Intn(2)) * time.Second
	time.Sleep(delay)

	// 90% success, 10% failure
	success := rand.Intn(10) != 0

	if success {
		payment.PaymentStatus = "succeeded"
		log.Println("Payment SUCCESS for order:", orderEvent.OrderId)
	} else {
		payment.PaymentStatus = "failed"
		log.Println("Payment FAILED for order:", orderEvent.OrderId)
	}

	err := payment.Save()
	if err != nil {
		log.Printf("Error saving the payment entry: %v", err)
	}

	//publish event to notify the payment status
	payload, err := json.Marshal(payment)
	if err != nil {
		log.Printf("Error marshaling payment for event: %v", err)
		return err
	}
	err = rabbitmq.PublishEvent("exchange", "topic", payload, "payment.status")
	if err != nil {
		return err
	}
	return nil
}
