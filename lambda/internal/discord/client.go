package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

// NewClient はDiscordのWebhookクライアントを作成する。
func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}

// PostWebhook はDiscordのWebhookにメッセージを送信する。
func (c *Client) PostWebhook(ctx context.Context, webhookURL string, msg *WebhookMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	return nil
}
