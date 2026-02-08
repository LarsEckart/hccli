package api

import "context"

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
	return List[Dataset](c, ctx, "/1/datasets")
}

func (c *Client) GetDataset(ctx context.Context, slug string) (*Dataset, error) {
	return Get[Dataset](c, ctx, "/1/datasets/"+slug)
}

func (c *Client) CreateDataset(ctx context.Context, ds *Dataset) (*Dataset, error) {
	return Create[Dataset](c, ctx, "/1/datasets", ds)
}

func (c *Client) UpdateDataset(ctx context.Context, slug string, ds *Dataset) (*Dataset, error) {
	return Update[Dataset](c, ctx, "/1/datasets/"+slug, ds)
}

func (c *Client) DeleteDataset(ctx context.Context, slug string) error {
	return Delete(c, ctx, "/1/datasets/"+slug)
}
