package models

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type MessagingManager struct {
	messagingClient *messaging.Client
}

func NewMessagingManager(app *firebase.App) *MessagingManager {
	messagingClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting messaging client: %v", err)
	}
	return &MessagingManager{
		messagingClient: messagingClient,
	}
}

func (mm *MessagingManager) SendNotification(token string) error {
	// Define the message payload
	message := &messaging.Message{
		Token: token, // Replace with the recipient's device token
		Notification: &messaging.Notification{
			Title: "Hello from Go!",
			Body:  "This is a test notification sent from Go",
		},
	}

	// Send the message
	response, err := mm.messagingClient.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	// Print the response
	fmt.Println("Successfully sent message:", response)
	return nil
}

func (mm *MessagingManager) SendDataMessage(token string) error {
	// Define the message payload
	message := &messaging.Message{
		Token: token, // Replace with the recipient's device token
		Data: map[string]string{
			"title":   "Test Tile",
			"message": "This is a test message from Go",
		},
	}

	// Send the message
	response, err := mm.messagingClient.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	// Print the response
	fmt.Println("Successfully sent message:", response)
	return nil
}
