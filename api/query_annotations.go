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
