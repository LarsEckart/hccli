package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type QueryAnnotation struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	QueryID     string `json:"query_id"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func (c *Client) CreateQueryAnnotation(ctx context.Context, dataset string, annotation *QueryAnnotation) (*QueryAnnotation, error) {
	body, err := json.Marshal(annotation)
	if err != nil {
		return nil, fmt.Errorf("encoding query annotation: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/query_annotations/"+dataset, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created QueryAnnotation
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) ListQueryAnnotations(ctx context.Context, dataset string) ([]QueryAnnotation, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/query_annotations/"+dataset, nil)
	if err != nil {
		return nil, err
	}

	var annotations []QueryAnnotation
	if err := c.do(req, &annotations); err != nil {
		return nil, err
	}
	return annotations, nil
}

func (c *Client) GetQueryAnnotation(ctx context.Context, dataset string, id string) (*QueryAnnotation, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/query_annotations/"+dataset+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	var annotation QueryAnnotation
	if err := c.do(req, &annotation); err != nil {
		return nil, err
	}
	return &annotation, nil
}

func (c *Client) UpdateQueryAnnotation(ctx context.Context, dataset string, id string, annotation *QueryAnnotation) (*QueryAnnotation, error) {
	body, err := json.Marshal(annotation)
	if err != nil {
		return nil, fmt.Errorf("encoding query annotation: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+"/1/query_annotations/"+dataset+"/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var updated QueryAnnotation
	if err := c.doJSON(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteQueryAnnotation(ctx context.Context, dataset string, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/1/query_annotations/"+dataset+"/"+id, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
