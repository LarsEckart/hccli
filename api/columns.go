package api

import "context"

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
	return List[Column](c, ctx, "/1/columns/"+dataset)
}

func (c *Client) GetColumn(ctx context.Context, dataset string, id string) (*Column, error) {
	return Get[Column](c, ctx, "/1/columns/"+dataset+"/"+id)
}

func (c *Client) CreateColumn(ctx context.Context, dataset string, col *Column) (*Column, error) {
	return Create[Column](c, ctx, "/1/columns/"+dataset, col)
}

func (c *Client) UpdateColumn(ctx context.Context, dataset string, id string, col *Column) (*Column, error) {
	return Update[Column](c, ctx, "/1/columns/"+dataset+"/"+id, col)
}

func (c *Client) DeleteColumn(ctx context.Context, dataset string, id string) error {
	return Delete(c, ctx, "/1/columns/"+dataset+"/"+id)
}
