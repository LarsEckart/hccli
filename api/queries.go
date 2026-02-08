package api

import "context"

type Query struct {
	ID                string        `json:"id,omitempty"`
	Breakdowns        []string      `json:"breakdowns,omitempty"`
	Calculations      []Calculation `json:"calculations,omitempty"`
	Filters           []QueryFilter `json:"filters,omitempty"`
	FilterCombination string        `json:"filter_combination,omitempty"`
	Granularity       int           `json:"granularity,omitempty"`
	Orders            []Order       `json:"orders,omitempty"`
	Limit             int           `json:"limit,omitempty"`
	StartTime         int           `json:"start_time,omitempty"`
	EndTime           int           `json:"end_time,omitempty"`
	TimeRange         int           `json:"time_range,omitempty"`
	Havings           []Having      `json:"havings,omitempty"`
}

type Calculation struct {
	Op     string `json:"op"`
	Column string `json:"column,omitempty"`
	Name   string `json:"name,omitempty"`
}

type QueryFilter struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Value  any    `json:"value,omitempty"`
}

type Order struct {
	Column string `json:"column,omitempty"`
	Op     string `json:"op,omitempty"`
	Order  string `json:"order,omitempty"`
}

type Having struct {
	CalculateOp string `json:"calculate_op"`
	Column      string `json:"column,omitempty"`
	Op          string `json:"op"`
	Value       any    `json:"value"`
}

func (c *Client) GetQuery(ctx context.Context, dataset string, queryID string) (*Query, error) {
	return Get[Query](c, ctx, "/1/queries/"+dataset+"/"+queryID)
}

func (c *Client) CreateQuery(ctx context.Context, dataset string, query *Query) (*Query, error) {
	return Create[Query](c, ctx, "/1/queries/"+dataset, query)
}
