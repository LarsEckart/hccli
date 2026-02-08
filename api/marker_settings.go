package api

import "context"

type MarkerSetting struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (c *Client) ListMarkerSettings(ctx context.Context, dataset string) ([]MarkerSetting, error) {
	return List[MarkerSetting](c, ctx, "/1/marker_settings/"+dataset)
}

func (c *Client) CreateMarkerSetting(ctx context.Context, dataset string, ms *MarkerSetting) (*MarkerSetting, error) {
	return Create[MarkerSetting](c, ctx, "/1/marker_settings/"+dataset, ms)
}

func (c *Client) UpdateMarkerSetting(ctx context.Context, dataset string, id string, ms *MarkerSetting) (*MarkerSetting, error) {
	return Update[MarkerSetting](c, ctx, "/1/marker_settings/"+dataset+"/"+id, ms)
}

func (c *Client) DeleteMarkerSetting(ctx context.Context, dataset string, id string) error {
	return Delete(c, ctx, "/1/marker_settings/"+dataset+"/"+id)
}
