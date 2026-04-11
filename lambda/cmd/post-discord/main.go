package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/caarlos0/env/v11"
	"github.com/ok-yyyy/discord-messenger/lambda/internal/aws/ssm"
	"github.com/ok-yyyy/discord-messenger/lambda/internal/discord"
	"github.com/ok-yyyy/discord-messenger/lambda/internal/message"
)

type Config struct {
	DiscordWebhookParameterName string `env:"DISCORD_WEBHOOK_PARAMETER_NAME,required"`
}

var (
	ssmClient     *ssm.Client
	discordClient *discord.Client
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("failed to load AWS config", slog.Any("error", err))
		os.Exit(1)
	}
	ssmClient = ssm.NewClient(cfg)
	discordClient = discord.NewClient()
}

func newConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return &cfg, nil
}

func handler(ctx context.Context, event events.SNSEvent) error {
	cfg, err := newConfig()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return err
	}

	// webhook urlをSSMパラメータストアから取得
	webhookURL, err := ssmClient.GetValue(ctx, cfg.DiscordWebhookParameterName, true)
	if err != nil {
		slog.Error("failed to get Discord webhook URL", slog.Any("error", err))
		return err
	}

	// SNSイベントを処理
	for _, record := range event.Records {
		slog.Info("processing SNS message", slog.String("messageId", record.SNS.MessageID))

		message, err := message.Unmarshal([]byte(record.SNS.Message))
		if err != nil {
			slog.Error("failed to unmarshal sns message", slog.Any("error", err))
			return err
		}

		webhookMessage, err := discord.BuildWebhookMessage(message)
		if err != nil {
			slog.Error("failed to build discord webhook message", slog.Any("error", err))
			return err
		}

		// メッセージ送信
		if err := discordClient.PostWebhook(ctx, webhookURL, webhookMessage); err != nil {
			slog.Error("failed to post webhook message", slog.Any("error", err))
			return err
		}
	}

	slog.Info("successfully processed SNS event")
	return nil
}

func main() {
	lambda.Start(handler)
}
