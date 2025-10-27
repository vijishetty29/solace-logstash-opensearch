package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Azure/go-amqp" // This is the AMQP 1.0 library
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func forwardToLogstash(host string, body []byte) {
	conn, err := net.DialTimeout("tcp", host, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to Logstash: %s", err)
		return
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\n", string(body))
	if err != nil {
		log.Printf("Failed to send data to Logstash: %s", err)
		return
	}
	log.Printf("Successfully forwarded message to Logstash")
}

func main() {
	solaceHost := getEnv("SOLACE_HOST", "test@localhost")
	// solaceUser := getEnv("SOLACE_USER", "admin")
	// solacePass := getEnv("SOLACE_PASS", "admin")
	logstashHost := getEnv("LOGSTASH_TCP_HOST", "localhost:5044")
	queueName := getEnv("SOLACE_QUEUE", "logstash_ingest_queue")
	solaceURL := fmt.Sprintf("amqp://%s:5672", solaceHost)

	log.Println("--- STARTING CONNECTOR ---")
	log.Printf("Connecting to Solace Host: %s (Full URL: %s)", solaceHost, solaceURL)
	//log.Printf("Using Solace User: %s", solaceUser)
	log.Println("--------------------------")

	ctx := context.Background()

	opts := &amqp.ConnOptions{
		//SASLType: amqp.SASLTypePlain(solaceUser, solacePass),
	}

	// Try to connect ONE TIME.
	// The 'depends_on' in docker-compose ensures Solace is healthy.
	client, err := amqp.Dial(ctx, solaceURL, opts)
	if err != nil {
		log.Fatalf("Failed to connect to Solace: %s", err)
	}
	defer client.Close()

	log.Println("Connected to Solace via AMQP 1.0")

	session, err := client.NewSession(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to open session: %s", err)
	}

	receiver, err := session.NewReceiver(ctx, queueName, nil)
	if err != nil {
		log.Fatalf("Failed to create receiver for queue '%s': %s", queueName, err)
	}
	defer receiver.Close(ctx)

	log.Printf("Receiver created. Waiting for messages on queue '%s'...", queueName)

	for {
		msg, err := receiver.Receive(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to receive message: %s", err)
		}

		log.Printf("Received message: %s", msg.GetData())
		forwardToLogstash(logstashHost, msg.GetData())

		err = receiver.AcceptMessage(ctx, msg)
		if err != nil {
			log.Printf("Warning: Failed to accept message: %s", err)
		}
	}
}
