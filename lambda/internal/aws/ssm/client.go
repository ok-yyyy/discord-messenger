package ssm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Client struct {
	client *ssm.Client
}

// NewClient はSSMクライアントを作成する。
func NewClient(cfg aws.Config) *Client {
	return &Client{
		client: ssm.NewFromConfig(cfg),
	}
}

// GetValue は指定されたパラメータの値を取得する。
func (c *Client) GetValue(ctx context.Context, name string, isSecure bool) (string, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(isSecure),
	}

	out, err := c.client.GetParameter(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get parameter %s: %w", name, err)
	}

	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", fmt.Errorf("parameter %s returned empty value", name)
	}

	return *out.Parameter.Value, nil
}
