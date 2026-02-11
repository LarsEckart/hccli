package api

import (
	"context"
	"fmt"
	"net/http"
)

type BurnAlert struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	AlertType   string `json:"alert_type,omitempty"`
	Triggered   *bool  `json:"triggered,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`

	// exhaustion_time fields
	ExhaustionMinutes *int `json:"exhaustion_minutes,omitempty"`

	// budget_rate fields
	BudgetRateWindowMinutes               *int `json:"budget_rate_window_minutes,omitempty"`
	BudgetRateDecreaseThresholdPerMillion *int `json:"budget_rate_decrease_threshold_per_million,omitempty"`

	SLO        *BurnAlertSLO           `json:"slo,omitempty"`
	Recipients []NotificationRecipient `json:"recipients,omitempty"`
}

type BurnAlertSLO struct {
	ID string `json:"id"`
}

type NotificationRecipient struct {
	ID      string                        `json:"id,omitempty"`
	Type    string                        `json:"type,omitempty"`
	Target  string                        `json:"target,omitempty"`
	Details *NotificationRecipientDetails `json:"details,omitempty"`
}

type NotificationRecipientDetails struct {
	PagerDutySeverity string            `json:"pagerduty_severity,omitempty"`
	Variables         []WebhookVariable `json:"variables,omitempty"`
}

type WebhookVariable struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

func (c *Client) ListBurnAlerts(ctx context.Context, dataset string, sloID string) ([]BurnAlert, error) {
	url := fmt.Sprintf("%s/1/burn_alerts/%s?slo_id=%s", c.BaseURL, dataset, sloID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result []BurnAlert
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetBurnAlert(ctx context.Context, dataset string, id string) (*BurnAlert, error) {
	return Get[BurnAlert](c, ctx, "/1/burn_alerts/"+dataset+"/"+id)
}

func (c *Client) CreateBurnAlert(ctx context.Context, dataset string, ba *BurnAlert) (*BurnAlert, error) {
	return Create[BurnAlert](c, ctx, "/1/burn_alerts/"+dataset, ba)
}

func (c *Client) UpdateBurnAlert(ctx context.Context, dataset string, id string, ba *BurnAlert) (*BurnAlert, error) {
	return Update[BurnAlert](c, ctx, "/1/burn_alerts/"+dataset+"/"+id, ba)
}

func (c *Client) DeleteBurnAlert(ctx context.Context, dataset string, id string) error {
	return Delete(c, ctx, "/1/burn_alerts/"+dataset+"/"+id)
}
