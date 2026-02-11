package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type QueryResult struct {
	ID       string         `json:"id,omitempty"`
	Complete bool           `json:"complete"`
	QueryID  string         `json:"query_id,omitempty"`
	Links    map[string]any `json:"links,omitempty"`
	Data     QueryData      `json:"data,omitempty"`
}

type QueryData struct {
	Series  []map[string]any `json:"series,omitempty"`
	Results []map[string]any `json:"results,omitempty"`
}

type createQueryResultRequest struct {
	QueryID string `json:"query_id"`
}

func (c *Client) CreateQueryResult(ctx context.Context, dataset string, queryID string) (*QueryResult, error) {
	body, err := json.Marshal(createQueryResultRequest{QueryID: queryID})
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/query_results/"+dataset, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result QueryResult
	if err := c.doJSON(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetQueryResult(ctx context.Context, dataset string, resultID string) (*QueryResult, error) {
	return Get[QueryResult](c, ctx, "/1/query_results/"+dataset+"/"+resultID)
}
