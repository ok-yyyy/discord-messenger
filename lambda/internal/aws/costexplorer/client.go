package costexplorer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

const (
	metricUnblendedCost = "UnblendedCost"
)

type Client struct {
	client *costexplorer.Client
}

type MonthlyCost struct {
	Period                 string
	Currency               string
	MonthToDateCost        string
	ForecastedMonthEndCost string
}

type costAmount struct {
	Currency string
	Amount   string
}

// NewClient はCost Explorerクライアントを作成する。
func NewClient(cfg aws.Config) *Client {
	return &Client{
		client: costexplorer.NewFromConfig(cfg),
	}
}

// GetMonthlyCost は指定された月の累計コストと予測コストを取得する。
func (c *Client) GetMonthlyCost(ctx context.Context, now time.Time) (*MonthlyCost, error) {
	now = now.UTC()

	monthToDateCost, err := c.getMonthToDateCost(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get month-to-date cost: %w", err)
	}

	forecastedMonthEndCost, err := c.getForecastedMonthEndCost(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecasted month end cost: %w", err)
	}

	if monthToDateCost.Currency != forecastedMonthEndCost.Currency {
		return nil, fmt.Errorf(
			"currency mismatch: monthToDate=%s, forecastedMonthEnd=%s",
			monthToDateCost.Currency,
			forecastedMonthEndCost.Currency,
		)
	}

	return &MonthlyCost{
		Period:                 now.Format("2006-01"),
		Currency:               monthToDateCost.Currency,
		MonthToDateCost:        monthToDateCost.Amount,
		ForecastedMonthEndCost: forecastedMonthEndCost.Amount,
	}, nil
}

// getMonthToDateCost は指定された月の累計コストを取得する。
func (c *Client) getMonthToDateCost(ctx context.Context, now time.Time) (*costAmount, error) {
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	out, err := c.client.GetCostAndUsage(ctx, &costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityMonthly,
		Metrics: []string{
			metricUnblendedCost,
		},
		TimePeriod: &types.DateInterval{
			Start: aws.String(monthStart.Format("2006-01-02")),
			End:   aws.String(nextMonthStart.Format("2006-01-02")),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cost: %w", err)
	}
	if len(out.ResultsByTime) != 1 {
		return nil, fmt.Errorf("unexpected results count: %d", len(out.ResultsByTime))
	}

	result := out.ResultsByTime[0]
	metric, ok := result.Total[metricUnblendedCost]
	if !ok {
		return nil, fmt.Errorf("metric %s not found", metricUnblendedCost)
	}

	slog.Info(
		"get month-to-date cost",
		slog.String("currency", aws.ToString(metric.Unit)),
		slog.String("amount", aws.ToString(metric.Amount)),
	)
	return &costAmount{
		Currency: aws.ToString(metric.Unit),
		Amount:   aws.ToString(metric.Amount),
	}, nil
}

// getForecastedMonthEndCost は指定された月の予測コストを取得する。
func (c *Client) getForecastedMonthEndCost(ctx context.Context, now time.Time) (*costAmount, error) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	nextMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)

	out, err := c.client.GetCostForecast(ctx, &costexplorer.GetCostForecastInput{
		Granularity: types.GranularityMonthly,
		Metric:      types.MetricUnblendedCost,
		TimePeriod: &types.DateInterval{
			Start: aws.String(today.Format("2006-01-02")),
			End:   aws.String(nextMonthStart.Format("2006-01-02")),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cost: %w", err)
	}
	if out.Total == nil {
		return nil, fmt.Errorf("forecast total is nil")
	}

	slog.Info(
		"get forecasted month end cost",
		slog.String("currency", aws.ToString(out.Total.Unit)),
		slog.String("amount", aws.ToString(out.Total.Amount)),
	)
	return &costAmount{
		Currency: aws.ToString(out.Total.Unit),
		Amount:   aws.ToString(out.Total.Amount),
	}, nil
}
