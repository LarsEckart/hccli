package api

import "context"

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
	return List[Marker](c, ctx, "/1/markers/"+dataset)
}

func (c *Client) CreateMarker(ctx context.Context, dataset string, m *Marker) (*Marker, error) {
	return Create[Marker](c, ctx, "/1/markers/"+dataset, m)
}

func (c *Client) UpdateMarker(ctx context.Context, dataset string, id string, m *Marker) (*Marker, error) {
	return Update[Marker](c, ctx, "/1/markers/"+dataset+"/"+id, m)
}

func (c *Client) DeleteMarker(ctx context.Context, dataset string, id string) error {
	return Delete(c, ctx, "/1/markers/"+dataset+"/"+id)
}
