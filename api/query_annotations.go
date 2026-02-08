package api

import "context"

type QueryAnnotation struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	QueryID     string `json:"query_id"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func (c *Client) CreateQueryAnnotation(ctx context.Context, dataset string, annotation *QueryAnnotation) (*QueryAnnotation, error) {
	return Create[QueryAnnotation](c, ctx, "/1/query_annotations/"+dataset, annotation)
}

func (c *Client) ListQueryAnnotations(ctx context.Context, dataset string) ([]QueryAnnotation, error) {
	return List[QueryAnnotation](c, ctx, "/1/query_annotations/"+dataset)
}

func (c *Client) GetQueryAnnotation(ctx context.Context, dataset string, id string) (*QueryAnnotation, error) {
	return Get[QueryAnnotation](c, ctx, "/1/query_annotations/"+dataset+"/"+id)
}

func (c *Client) UpdateQueryAnnotation(ctx context.Context, dataset string, id string, annotation *QueryAnnotation) (*QueryAnnotation, error) {
	return Update[QueryAnnotation](c, ctx, "/1/query_annotations/"+dataset+"/"+id, annotation)
}

func (c *Client) DeleteQueryAnnotation(ctx context.Context, dataset string, id string) error {
	return Delete(c, ctx, "/1/query_annotations/"+dataset+"/"+id)
}
