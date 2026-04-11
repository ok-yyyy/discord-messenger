package sns

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Client struct {
	client *sns.Client
}

// NewClient はSNSクライアントを作成する。
func NewClient(cfg aws.Config) *Client {
	return &Client{
		client: sns.NewFromConfig(cfg),
	}
}

// Publish は指定されたTopicにメッセージを送信する。
func (c *Client) Publish(ctx context.Context, topicARN string, message string) error {
	input := &sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(message),
	}

	if _, err := c.client.Publish(ctx, input); err != nil {
		return fmt.Errorf("failed to publish message to topic %s: %w", topicARN, err)
	}

	return nil
}
