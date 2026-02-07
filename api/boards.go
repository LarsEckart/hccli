package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (c *Client) ListBoards(ctx context.Context) ([]Board, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/1/boards", nil)
	if err != nil {
		return nil, err
	}

	var boards []Board
	if err := c.do(req, &boards); err != nil {
		return nil, err
	}
	return boards, nil
}

func (c *Client) DeleteBoard(ctx context.Context, boardID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/1/boards/"+boardID, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) CreateBoard(ctx context.Context, board *Board) (*Board, error) {
	body, err := json.Marshal(board)
	if err != nil {
		return nil, fmt.Errorf("encoding board: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/boards", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created Board
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) CreateBoardView(ctx context.Context, boardID string, view *BoardView) (*BoardView, error) {
	body, err := json.Marshal(view)
	if err != nil {
		return nil, fmt.Errorf("encoding board view: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/1/boards/"+boardID+"/views", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var created BoardView
	if err := c.doJSON(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}
