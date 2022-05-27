package push

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Client interface {
	PushNotification(ctx context.Context, token, title, body string, payload map[string]string) error

	GoatPostPushNotification(ctx context.Context, postID, goatID int, title, body string) error
	SubscribeToGoatPostTopic(ctx context.Context, goatID int, token string) error
	UnsubscribeFromGoatPostTopic(ctx context.Context, goatID int, token string) error
}

type client struct {
	fcmClient *messaging.Client
}

func NewClient(ctx context.Context) (Client, error) {
	c := client{}
	opts := []option.ClientOption{}
	app, newAppErr := firebase.NewApp(ctx, nil, opts...)
	if newAppErr != nil {
		return nil, newAppErr
	}

	fcmClient, fcmClientErr := app.Messaging(ctx)
	if fcmClientErr != nil {
		return nil, fcmClientErr
	}
	c.fcmClient = fcmClient
	return &c, nil
}

func (c *client) PushNotification(ctx context.Context, token, title, body string, payload map[string]string) error {
	// https://firebase.google.com/docs/cloud-messaging/send-message
	_, pushNotificationErr := c.fcmClient.Send(ctx, &messaging.Message{
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					MutableContent: true,
				},
			},
		},
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
		Data:  payload,
	})
	return pushNotificationErr
}

const goatPostTopicString = "goat-%d-posts" // where %d is the creator user ID
const goatPostNotificationType = "POST"

func (c *client) GoatPostPushNotification(ctx context.Context, postID, goatID int, title, body string) error {
	payloadData := map[string]string{
		"type": goatPostNotificationType,
		"id":   fmt.Sprint(postID),
	}
	topic := fmt.Sprintf(goatPostTopicString, goatID)
	// https://firebase.google.com/docs/cloud-messaging/send-message
	_, pushNotificationErr := c.fcmClient.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Topic: topic,
		Data:  payloadData,
	})
	return pushNotificationErr
}

func (c *client) SubscribeToGoatPostTopic(ctx context.Context, goatID int, token string) error {
	topic := fmt.Sprintf(goatPostTopicString, goatID)
	tokens := []string{token}
	_, topicSubscribeErr := c.fcmClient.SubscribeToTopic(ctx, tokens, topic)
	return topicSubscribeErr
}

func (c *client) UnsubscribeFromGoatPostTopic(ctx context.Context, goatID int, token string) error {
	topic := fmt.Sprintf(goatPostTopicString, goatID)
	tokens := []string{token}
	_, topicSubscribeErr := c.fcmClient.UnsubscribeFromTopic(ctx, tokens, topic)
	return topicSubscribeErr
}
