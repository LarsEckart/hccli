package api

import "context"

type Board struct {
	ID               string         `json:"id,omitempty"`
	Name             string         `json:"name"`
	Description      string         `json:"description,omitempty"`
	Type             string         `json:"type"`
	Links            *BoardLinks    `json:"links,omitempty"`
	Panels           []BoardPanel   `json:"panels,omitempty"`
	Tags             []Tag          `json:"tags,omitempty"`
	PresetFilters    []PresetFilter `json:"preset_filters,omitempty"`
	LayoutGeneration string         `json:"layout_generation,omitempty"`
}

type BoardLinks struct {
	BoardURL string `json:"board_url"`
}

type BoardPanel struct {
	Type       string      `json:"type"`
	QueryPanel *QueryPanel `json:"query_panel,omitempty"`
	SLOPanel   *SLOPanel   `json:"slo_panel,omitempty"`
	TextPanel  *TextPanel  `json:"text_panel,omitempty"`
	Position   *Position   `json:"position,omitempty"`
}

type QueryPanel struct {
	QueryID           string `json:"query_id"`
	QueryAnnotationID string `json:"query_annotation_id,omitempty"`
	QueryStyle        string `json:"query_style,omitempty"`
}

type SLOPanel struct {
	SLOID string `json:"slo_id"`
}

type TextPanel struct {
	Content string `json:"content"`
}

type Position struct {
	XCoordinate int `json:"x_coordinate"`
	YCoordinate int `json:"y_coordinate"`
	Height      int `json:"height"`
	Width       int `json:"width"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PresetFilter struct {
	Column string `json:"column"`
	Alias  string `json:"alias,omitempty"`
}

type BoardView struct {
	ID      string            `json:"id,omitempty"`
	Name    string            `json:"name"`
	Filters []BoardViewFilter `json:"filters"`
}

type BoardViewFilter struct {
	Column    string `json:"column"`
	Operation string `json:"operation"`
	Value     any    `json:"value,omitempty"`
}

func (c *Client) GetBoard(ctx context.Context, boardID string) (*Board, error) {
	return Get[Board](c, ctx, "/1/boards/"+boardID)
}

func (c *Client) ListBoards(ctx context.Context) ([]Board, error) {
	return List[Board](c, ctx, "/1/boards")
}

func (c *Client) DeleteBoard(ctx context.Context, boardID string) error {
	return Delete(c, ctx, "/1/boards/"+boardID)
}

func (c *Client) CreateBoard(ctx context.Context, board *Board) (*Board, error) {
	return Create[Board](c, ctx, "/1/boards", board)
}

func (c *Client) UpdateBoard(ctx context.Context, boardID string, board *Board) (*Board, error) {
	return Update[Board](c, ctx, "/1/boards/"+boardID, board)
}

func (c *Client) ListBoardViews(ctx context.Context, boardID string) ([]BoardView, error) {
	return List[BoardView](c, ctx, "/1/boards/"+boardID+"/views")
}

func (c *Client) GetBoardView(ctx context.Context, boardID, viewID string) (*BoardView, error) {
	return Get[BoardView](c, ctx, "/1/boards/"+boardID+"/views/"+viewID)
}

func (c *Client) CreateBoardView(ctx context.Context, boardID string, view *BoardView) (*BoardView, error) {
	return Create[BoardView](c, ctx, "/1/boards/"+boardID+"/views", view)
}

func (c *Client) UpdateBoardView(ctx context.Context, boardID, viewID string, view *BoardView) (*BoardView, error) {
	return Update[BoardView](c, ctx, "/1/boards/"+boardID+"/views/"+viewID, view)
}

func (c *Client) DeleteBoardView(ctx context.Context, boardID, viewID string) error {
	return Delete(c, ctx, "/1/boards/"+boardID+"/views/"+viewID)
}
