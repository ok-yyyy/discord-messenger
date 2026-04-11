package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/ok-yyyy/discord-messenger/lambda/internal/message"
)

const (
	colorGreen = 0x2ecc71
)

// https://docs.discord.com/developers/resources/webhook#execute-webhook
type WebhookMessage struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Embeds    []Embed `json:"embeds"`
}

// https://docs.discord.com/developers/resources/message#embed-object
type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Timestamp   string       `json:"timestamp,omitempty"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// BuildWebhookMessage はwebhook用のメッセージを作成する。
func BuildWebhookMessage(msg message.Message) (*WebhookMessage, error) {
	switch msg.Type() {
	case message.TypeSimple:
		return buildSimpleMessage(msg.(*message.SimpleMessage)), nil

	case message.TypeAWSCost:
		return buildAWSCostMessage(msg.(*message.AWSCostMessage)), nil

	default:
		return nil, fmt.Errorf("unknown message type: %s", msg.Type())
	}
}

func buildSimpleMessage(msg *message.SimpleMessage) *WebhookMessage {
	return &WebhookMessage{
		Embeds: []Embed{
			{
				Title:       msg.Title,
				Description: msg.Description,
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
				Color:       colorGreen,
			},
		},
	}
}

func buildAWSCostMessage(msg *message.AWSCostMessage) *WebhookMessage {
	truncate := func(s string) string {
		whole, fraction, found := strings.Cut(s, ".")
		if !found {
			return whole + ".00"
		}
		if len(fraction) >= 2 {
			return whole + "." + fraction[:2]
		}
		if len(fraction) == 1 {
			return whole + "." + fraction + "0"
		}
		return whole + ".00"
	}

	monthToDateCost := truncate(msg.MonthToDateCost)
	forecastedMonthEndCost := truncate(msg.ForecastedMonthEndCost)

	return &WebhookMessage{
		Username:  "AWS Cost",
		AvatarURL: "https://github.com/aws.png",
		Embeds: []Embed{
			{
				Title:       "AWS Cost Notification",
				Description: msg.Period,
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
				Color:       colorGreen,
				Fields: []EmbedField{
					{Name: "累計コスト", Value: fmt.Sprintf("%s %s", msg.Currency, monthToDateCost)},
					{Name: "予測コスト", Value: fmt.Sprintf("%s %s", msg.Currency, forecastedMonthEndCost)},
				},
			},
		},
	}
}
