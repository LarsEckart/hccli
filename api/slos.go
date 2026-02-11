package api

import "context"

type SLO struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	SLI              SLOSLI   `json:"sli"`
	TimePeriodDays   int      `json:"time_period_days"`
	TargetPerMillion int      `json:"target_per_million"`
	Tags             []Tag    `json:"tags,omitempty"`
	DatasetSlugs     []string `json:"dataset_slugs,omitempty"`
	ResetAt          *string  `json:"reset_at,omitempty"`
	CreatedAt        string   `json:"created_at,omitempty"`
	UpdatedAt        string   `json:"updated_at,omitempty"`

	// Detailed response fields (returned when ?detailed is set)
	Compliance      *float64 `json:"compliance,omitempty"`
	BudgetRemaining *float64 `json:"budget_remaining,omitempty"`
	Status          string   `json:"status,omitempty"`
	BurnRate        *float64 `json:"burn_rate,omitempty"`
}

type SLOSLI struct {
	Alias string `json:"alias"`
}

func (c *Client) ListSLOs(ctx context.Context, dataset string) ([]SLO, error) {
	return List[SLO](c, ctx, "/1/slos/"+dataset)
}

func (c *Client) GetSLO(ctx context.Context, dataset string, id string) (*SLO, error) {
	return Get[SLO](c, ctx, "/1/slos/"+dataset+"/"+id)
}

func (c *Client) GetSLODetailed(ctx context.Context, dataset string, id string) (*SLO, error) {
	return Get[SLO](c, ctx, "/1/slos/"+dataset+"/"+id+"?detailed")
}

func (c *Client) CreateSLO(ctx context.Context, dataset string, slo *SLO) (*SLO, error) {
	return Create[SLO](c, ctx, "/1/slos/"+dataset, slo)
}

func (c *Client) UpdateSLO(ctx context.Context, dataset string, id string, slo *SLO) (*SLO, error) {
	return Update[SLO](c, ctx, "/1/slos/"+dataset+"/"+id, slo)
}

func (c *Client) DeleteSLO(ctx context.Context, dataset string, id string) error {
	return Delete(c, ctx, "/1/slos/"+dataset+"/"+id)
}
