package message

import "encoding/json"

type Type string

const (
	TypeSimple  Type = "simple"
	TypeAWSCost Type = "aws_cost"
)

type Message interface {
	Type() Type
}

type Envelope struct {
	Type Type            `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Simple Message

type SimpleMessage struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (m *SimpleMessage) Type() Type {
	return TypeSimple
}

// AWS Cost Message

type AWSCostMessage struct {
	Period                 string `json:"period"`
	Currency               string `json:"currency"`
	MonthToDateCost        string `json:"month_to_date_cost"`
	ForecastedMonthEndCost string `json:"forecasted_month_end_cost"`
}

func (m *AWSCostMessage) Type() Type {
	return TypeAWSCost
}

// Ensure all message types implement the Message interface

var _ Message = (*SimpleMessage)(nil)
var _ Message = (*AWSCostMessage)(nil)
