package api

import (
	"context"
	"net/http"
)

type AuthResponse struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	APIKeyAccess APIKeyAccess `json:"api_key_access"`
	Environment  Environment  `json:"environment"`
	Team         Team         `json:"team"`
}

type APIKeyAccess struct {
	Events         bool `json:"events"`
	Markers        bool `json:"markers"`
	Triggers       bool `json:"triggers"`
	Boards         bool `json:"boards"`
	Queries        bool `json:"queries"`
	Columns        bool `json:"columns"`
	CreateDatasets bool `json:"createDatasets"`
	SLOs           bool `json:"slos"`
	Recipients     bool `json:"recipients"`
	PrivateBoards  bool `json:"privateBoards"`
}

type Environment struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Team struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (c *Client) GetAuth(ctx context.Context) (*AuthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/auth", nil)
	if err != nil {
		return nil, err
	}

	var auth AuthResponse
	if err := c.do(req, &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}

type AuthV2Response struct {
	Included []AuthV2IncludedResource `json:"included"`
	Data     AuthV2Data               `json:"data"`
}

type AuthV2IncludedResource struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type AuthV2Data struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    AuthV2Attributes       `json:"attributes"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

type AuthV2Attributes struct {
	Name       string            `json:"name"`
	KeyType    string            `json:"key_type"`
	Disabled   bool              `json:"disabled"`
	Scopes     []string          `json:"scopes"`
	Timestamps *AuthV2Timestamps `json:"timestamps,omitempty"`
}

type AuthV2Timestamps struct {
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (c *Client) GetAuthV2(ctx context.Context) (*AuthV2Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/2/auth", nil)
	if err != nil {
		return nil, err
	}

	var auth AuthV2Response
	if err := c.doBearer(req, &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}
