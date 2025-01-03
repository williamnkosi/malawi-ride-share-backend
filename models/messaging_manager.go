package models

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
)

type MessagingManager struct {
}

func NewMessagingManager() *MessagingManager {
	return &MessagingManager{}
}

func (mm *MessagingManager) SendNotification(client *messaging.Client, token string) error {
	// Define the message payload
	message := &messaging.Message{
		Token: token, // Replace with the recipient's device token
		Notification: &messaging.Notification{
			Title: "Hello from Go!",
			Body:  "This is a test notification sent from Go",
		},
	}

	// Send the message
	response, err := client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	// Print the response
	fmt.Println("Successfully sent message:", response)
	return nil
}
