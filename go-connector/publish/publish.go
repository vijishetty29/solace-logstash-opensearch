package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-amqp"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	// This field will be omitted because it's not exported
	internalID string
}

func main() {
	// 1. Set your Solace broker's AMQP connection details
	// Make sure the AMQP service is enabled on your Solace Message VPN
	host := "localhost" // Or "localhost"
	port := 5672        // Default AMQPS port

	// Create the connection string
	// "amqps" = AMQP over SSL/TLS
	// "amqp" = AMQP (plaintext, usually port 5672)
	connStr := fmt.Sprintf("amqp://%s:%d", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Connect to the broker
	client, err := amqp.Dial(ctx, connStr, nil)
	if err != nil {
		log.Fatalf("Failed to dial AMQP broker: %v", err)
	}
	defer client.Close()

	// 3. Open a session
	session, err := client.NewSession(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create AMQP session: %v", err)
	}

	// 4. Create a sender (publisher)
	// In Solace, the AMQP "target address" *is* the topic name.
	solaceTopic := "logstash_ingest_queue"
	sender, err := session.NewSender(ctx, solaceTopic, nil)
	if err != nil {
		log.Fatalf("Failed to create sender: %v", err)
	}
	defer sender.Close(ctx)
	user := User{
		Name:       "Alice",
		Age:        30,
		Email:      "alice@example.com",
		internalID: "xyz-123",
	}

	// 2. Marshal the struct into a byte slice
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// 3. Convert the byte slice to a string
	jsonString := string(jsonData)
	// 5. Create and send the message
	//messageBody := "Hello Solace from Go (via AMQP 1.0)!"
	msg := amqp.NewMessage([]byte(jsonString))

	log.Printf("Publishing to topic '%s': %s", solaceTopic, jsonString)
	err = sender.Send(ctx, msg, nil)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message published successfully!")
}
