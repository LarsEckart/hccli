package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DerivedColumn struct {
	ID          string `json:"id,omitempty"`
	Alias       string `json:"alias"`
	Expression  string `json:"expression"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func (c *Client) ListDerivedColumns(ctx context.Context, dataset string) ([]DerivedColumn, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/derived_columns/"+dataset, nil)
	if err != nil {
		return nil, err
	}

	var cols []DerivedColumn
	if err := c.do(req, &cols); err != nil {
		return nil, err
	}
	return cols, nil
}

func (c *Client) GetDerivedColumn(ctx context.Context, dataset string, id string) (*DerivedColumn, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/derived_columns/"+dataset+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	var col DerivedColumn
	if err := c.do(req, &col); err != nil {
		return nil, err
	}
	return &col, nil
}

func (c *Client) CreateDerivedColumn(ctx context.Context, dataset string, col *DerivedColumn) (*DerivedColumn, error) {
	body, err := json.Marshal(col)
	if err != nil {
		return nil, fmt.Errorf("encoding derived column: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/derived_columns/"+dataset, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created DerivedColumn
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateDerivedColumn(ctx context.Context, dataset string, id string, col *DerivedColumn) (*DerivedColumn, error) {
	body, err := json.Marshal(col)
	if err != nil {
		return nil, fmt.Errorf("encoding derived column: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+"/1/derived_columns/"+dataset+"/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var updated DerivedColumn
	if err := c.doJSON(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteDerivedColumn(ctx context.Context, dataset string, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/1/derived_columns/"+dataset+"/"+id, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
