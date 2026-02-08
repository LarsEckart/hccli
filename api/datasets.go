package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DatasetSettings struct {
	DeleteProtected *bool `json:"delete_protected,omitempty"`
}

type Dataset struct {
	Name                string           `json:"name,omitempty"`
	Description         string           `json:"description,omitempty"`
	Settings            *DatasetSettings `json:"settings,omitempty"`
	ExpandJSONDepth     *int             `json:"expand_json_depth,omitempty"`
	Slug                string           `json:"slug,omitempty"`
	RegularColumnsCount *int             `json:"regular_columns_count,omitempty"`
	LastWrittenAt       *string          `json:"last_written_at,omitempty"`
	CreatedAt           string           `json:"created_at,omitempty"`
}

func (c *Client) ListDatasets(ctx context.Context) ([]Dataset, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/datasets", nil)
	if err != nil {
		return nil, err
	}

	var datasets []Dataset
	if err := c.do(req, &datasets); err != nil {
		return nil, err
	}
	return datasets, nil
}

func (c *Client) GetDataset(ctx context.Context, slug string) (*Dataset, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/datasets/"+slug, nil)
	if err != nil {
		return nil, err
	}

	var ds Dataset
	if err := c.do(req, &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func (c *Client) CreateDataset(ctx context.Context, ds *Dataset) (*Dataset, error) {
	body, err := json.Marshal(ds)
	if err != nil {
		return nil, fmt.Errorf("encoding dataset: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/datasets", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created Dataset
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateDataset(ctx context.Context, slug string, ds *Dataset) (*Dataset, error) {
	body, err := json.Marshal(ds)
	if err != nil {
		return nil, fmt.Errorf("encoding dataset: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+"/1/datasets/"+slug, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var updated Dataset
	if err := c.doJSON(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteDataset(ctx context.Context, slug string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/1/datasets/"+slug, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
