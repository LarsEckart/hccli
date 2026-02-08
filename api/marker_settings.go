package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type MarkerSetting struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (c *Client) ListMarkerSettings(ctx context.Context, dataset string) ([]MarkerSetting, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/marker_settings/"+dataset, nil)
	if err != nil {
		return nil, err
	}

	var settings []MarkerSetting
	if err := c.do(req, &settings); err != nil {
		return nil, err
	}
	return settings, nil
}

func (c *Client) CreateMarkerSetting(ctx context.Context, dataset string, ms *MarkerSetting) (*MarkerSetting, error) {
	body, err := json.Marshal(ms)
	if err != nil {
		return nil, fmt.Errorf("encoding marker setting: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/marker_settings/"+dataset, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created MarkerSetting
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateMarkerSetting(ctx context.Context, dataset string, id string, ms *MarkerSetting) (*MarkerSetting, error) {
	body, err := json.Marshal(ms)
	if err != nil {
		return nil, fmt.Errorf("encoding marker setting: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+"/1/marker_settings/"+dataset+"/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var updated MarkerSetting
	if err := c.doJSON(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteMarkerSetting(ctx context.Context, dataset string, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/1/marker_settings/"+dataset+"/"+id, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
