package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer redisClient.Close()

	// Subscribe to Redis Channel
	channel := "task_event_channel"
	pubsub := redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	fmt.Printf("Subscribed to channel: %s\n", channel)

	// Listen for Events
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
			break
		}

		fmt.Printf("Received event: %s\n", msg.Payload)

		// Trigger Task Processing
		processTaskFromQueue(redisClient)
	}
}

func processTaskFromQueue(redisClient *redis.Client) {
	queueName := "task_queue"

	for {
		// Pop Task from Queue
		task, err := redisClient.LPop(ctx, queueName).Result()
		if err == redis.Nil {
			fmt.Println("No tasks in queue")
			break
		} else if err != nil {
			log.Printf("Error fetching task: %v", err)
			break
		}

		// Process Task
		fmt.Printf("Processing task: %s\n", task)
		time.Sleep(2 * time.Second) // Simulate task processing

	}
}
