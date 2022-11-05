package _const

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type contextKey string

func (c contextKey) String() string {
	return "context key " + string(c)
}

var (
	contextKeyPubsubClient = contextKey("pubsub-client")
	contextKeyPubsubTopic  = contextKey("pubsub-topic")
)

func SetPubsubClient(ctx context.Context, pubsubClient *pubsub.Client) context.Context {
	return context.WithValue(ctx, contextKeyPubsubClient, pubsubClient)
}

func GetPubsubClient(ctx context.Context) *pubsub.Client {
	return ctx.Value(contextKeyPubsubClient).(*pubsub.Client)
}

func SetPubsubTopic(ctx context.Context, topic *pubsub.Topic) context.Context {
	return context.WithValue(ctx, contextKey(fmt.Sprintf("%s-%s", contextKeyPubsubTopic, topic.ID())), topic)
}

func GetPubsubTopic(ctx context.Context, topicID string) *pubsub.Topic {
	return ctx.Value(contextKey(fmt.Sprintf("%s-%s", contextKeyPubsubTopic, topicID))).(*pubsub.Topic)
}
