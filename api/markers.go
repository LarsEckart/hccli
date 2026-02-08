package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Marker struct {
	ID        string `json:"id,omitempty"`
	StartTime *int64 `json:"start_time,omitempty"`
	EndTime   *int64 `json:"end_time,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
	URL       string `json:"url,omitempty"`
	Color     string `json:"color,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (c *Client) ListMarkers(ctx context.Context, dataset string) ([]Marker, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/markers/"+dataset, nil)
	if err != nil {
		return nil, err
	}

	var markers []Marker
	if err := c.do(req, &markers); err != nil {
		return nil, err
	}
	return markers, nil
}

func (c *Client) CreateMarker(ctx context.Context, dataset string, m *Marker) (*Marker, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encoding marker: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/markers/"+dataset, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created Marker
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateMarker(ctx context.Context, dataset string, id string, m *Marker) (*Marker, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encoding marker: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+"/1/markers/"+dataset+"/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var updated Marker
	if err := c.doJSON(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteMarker(ctx context.Context, dataset string, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/1/markers/"+dataset+"/"+id, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
