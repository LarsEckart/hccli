package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Column struct {
	ID          string `json:"id,omitempty"`
	KeyName     string `json:"key_name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Hidden      *bool  `json:"hidden,omitempty"`
	LastWritten string `json:"last_written,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func (c *Client) ListColumns(ctx context.Context, dataset string) ([]Column, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/columns/"+dataset, nil)
	if err != nil {
		return nil, err
	}

	var cols []Column
	if err := c.do(req, &cols); err != nil {
		return nil, err
	}
	return cols, nil
}

func (c *Client) GetColumn(ctx context.Context, dataset string, id string) (*Column, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/columns/"+dataset+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	var col Column
	if err := c.do(req, &col); err != nil {
		return nil, err
	}
	return &col, nil
}

func (c *Client) CreateColumn(ctx context.Context, dataset string, col *Column) (*Column, error) {
	body, err := json.Marshal(col)
	if err != nil {
		return nil, fmt.Errorf("encoding column: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/columns/"+dataset, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created Column
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateColumn(ctx context.Context, dataset string, id string, col *Column) (*Column, error) {
	body, err := json.Marshal(col)
	if err != nil {
		return nil, fmt.Errorf("encoding column: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+"/1/columns/"+dataset+"/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var updated Column
	if err := c.doJSON(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteColumn(ctx context.Context, dataset string, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/1/columns/"+dataset+"/"+id, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
