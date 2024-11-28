package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
	app := fiber.New()

	// Redis setup
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Redis container host
	})

	// Endpoint to publish a message
	app.Get("/publish", func(c *fiber.Ctx) error {

		// Push Task to Redis Queue
		queueName := "task_queue"
		task := "task: " + time.Now().String()
		err := redisClient.RPush(context.Background(), queueName, task).Err()
		if err != nil {
			log.Fatalf("Failed to push task to queue: %v", err)
		}
		fmt.Printf("Task pushed to queue: %s\n", task)

		// Publish Event to Redis Pub/Sub Channel
		channel := "task_event_channel"
		err = redisClient.Publish(context.Background(), channel, "new_task").Err()
		if err != nil {
			log.Fatalf("Failed to publish event: %v", err)
		}
		fmt.Printf("Event published to channel: %s\n", channel)

		return c.JSON(fiber.Map{"status": "Message published"})
	})

	log.Println("Publisher service running on port 8000")
	log.Fatal(app.Listen(":8000"))
}
