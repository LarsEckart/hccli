package api_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/LarsEckart/hccli/api"
)

func TestQueryHavingOmitsEmptyCalculateOp(t *testing.T) {
	query := api.Query{
		Havings: []api.Having{
			{
				Column: "error_rate",
				Op:     ">",
				Value:  0.01,
			},
		},
	}

	data, err := json.Marshal(query)
	if err != nil {
		t.Fatalf("marshal query: %v", err)
	}

	if bytes.Contains(data, []byte(`"calculate_op"`)) {
		t.Fatalf("expected calculate_op to be omitted, got %s", data)
	}
}
