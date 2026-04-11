package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/caarlos0/env/v11"
	"github.com/ok-yyyy/discord-messenger/lambda/internal/aws/costexplorer"
	"github.com/ok-yyyy/discord-messenger/lambda/internal/aws/sns"
	"github.com/ok-yyyy/discord-messenger/lambda/internal/message"
)

type Config struct {
	SNSTopicARN string `env:"SNS_TOPIC_ARN,required"`
}

var (
	costexplorerClient *costexplorer.Client
	snsClient          *sns.Client
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("failed to load AWS config", slog.Any("error", err))
		os.Exit(1)
	}
	costexplorerClient = costexplorer.NewClient(cfg)
	snsClient = sns.NewClient(cfg)
}

func newConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return &cfg, nil
}

func handler(ctx context.Context, event events.EventBridgeEvent) error {
	cfg, err := newConfig()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return err
	}

	// コストを取得
	cost, err := costexplorerClient.GetMonthlyCost(ctx, time.Now())
	if err != nil {
		slog.Error("failed to get month-to-date cost", slog.Any("error", err))
		return err
	}

	msg := &message.AWSCostMessage{
		Period:                 cost.Period,
		Currency:               cost.Currency,
		MonthToDateCost:        cost.MonthToDateCost,
		ForecastedMonthEndCost: cost.ForecastedMonthEndCost,
	}

	body, err := message.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal message", slog.Any("error", err))
		return err
	}

	// SNSにコスト情報を送信
	if err := snsClient.Publish(ctx, cfg.SNSTopicARN, string(body)); err != nil {
		slog.Error("failed to publish message", slog.Any("error", err))
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
