package api

import "context"

type DerivedColumn struct {
	ID          string `json:"id,omitempty"`
	Alias       string `json:"alias"`
	Expression  string `json:"expression"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func (c *Client) ListDerivedColumns(ctx context.Context, dataset string) ([]DerivedColumn, error) {
	return List[DerivedColumn](c, ctx, "/1/derived_columns/"+dataset)
}

func (c *Client) GetDerivedColumn(ctx context.Context, dataset string, id string) (*DerivedColumn, error) {
	return Get[DerivedColumn](c, ctx, "/1/derived_columns/"+dataset+"/"+id)
}

func (c *Client) CreateDerivedColumn(ctx context.Context, dataset string, col *DerivedColumn) (*DerivedColumn, error) {
	return Create[DerivedColumn](c, ctx, "/1/derived_columns/"+dataset, col)
}

func (c *Client) UpdateDerivedColumn(ctx context.Context, dataset string, id string, col *DerivedColumn) (*DerivedColumn, error) {
	return Update[DerivedColumn](c, ctx, "/1/derived_columns/"+dataset+"/"+id, col)
}

func (c *Client) DeleteDerivedColumn(ctx context.Context, dataset string, id string) error {
	return Delete(c, ctx, "/1/derived_columns/"+dataset+"/"+id)
}
