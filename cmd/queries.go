package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/LarsEckart/hccli/api"
	"github.com/LarsEckart/hccli/timefmt"
	"github.com/urfave/cli/v3"
)

// noValueOps are filter operators that take no value argument.
var noValueOps = map[string]bool{
	"exists":         true,
	"does-not-exist": true,
}

var calculationOps = map[string]bool{
	"AVG":            true,
	"CONCURRENCY":    true,
	"COUNT":          true,
	"COUNT_DISTINCT": true,
	"HEATMAP":        true,
	"MAX":            true,
	"MIN":            true,
	"P001":           true,
	"P01":            true,
	"P05":            true,
	"P10":            true,
	"P20":            true,
	"P25":            true,
	"P50":            true,
	"P75":            true,
	"P80":            true,
	"P90":            true,
	"P95":            true,
	"P99":            true,
	"P999":           true,
	"RATE_AVG":       true,
	"RATE_MAX":       true,
	"RATE_SUM":       true,
	"SUM":            true,
}

var orderDirections = map[string]string{
	"asc":        "ascending",
	"ascending":  "ascending",
	"desc":       "descending",
	"descending": "descending",
}

var havingOps = map[string]bool{
	"=":  true,
	"!=": true,
	">":  true,
	">=": true,
	"<":  true,
	"<=": true,
}

var queryFlagNames = []string{
	"calculation-op",
	"calculation-column",
	"breakdown",
	"filter",
	"filter-combination",
	"order",
	"limit",
	"having",
	"time-range",
	"from",
	"to",
	"timezone",
}

type queryInput struct {
	Query   *api.Query
	RawJSON []byte
}

// parseFilter parses a filter string in the form "column op [value]".
// The column is the first whitespace-delimited token, the op is the second,
// and the optional value is everything after the op.
func parseFilter(s string) (api.QueryFilter, error) {
	// Split into at most 3 parts: column, op, value
	parts := strings.SplitN(strings.TrimSpace(s), " ", 3)
	if len(parts) < 2 {
		return api.QueryFilter{}, fmt.Errorf("invalid filter %q: expected \"column op [value]\"", s)
	}

	col := parts[0]
	op := parts[1]

	f := api.QueryFilter{
		Column: col,
		Op:     op,
	}

	if noValueOps[op] {
		return f, nil
	}

	if len(parts) < 3 || parts[2] == "" {
		return api.QueryFilter{}, fmt.Errorf("invalid filter %q: operator %q requires a value", s, op)
	}
	f.Value = parts[2]
	return f, nil
}

func parseOrder(s string) (api.Order, error) {
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) != 2 {
		return api.Order{}, fmt.Errorf("invalid order %q: expected \"term asc|desc\" (for example, \"MAX(duration_ms) desc\")", s)
	}

	direction, ok := orderDirections[strings.ToLower(parts[1])]
	if !ok {
		return api.Order{}, fmt.Errorf("invalid order %q: direction must be asc, ascending, desc, or descending", s)
	}

	op, column, err := parseQueryTerm(parts[0], true)
	if err != nil {
		return api.Order{}, fmt.Errorf("invalid order %q: %w", s, err)
	}

	return api.Order{
		Op:     op,
		Column: column,
		Order:  direction,
	}, nil
}

func parseHaving(s string) (api.Having, error) {
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) != 3 {
		return api.Having{}, fmt.Errorf("invalid having %q: expected \"CALC(column) op number\" (for example, \"MAX(duration_ms) > 1000\")", s)
	}

	if !havingOps[parts[1]] {
		return api.Having{}, fmt.Errorf("invalid having %q: operator must be one of =, !=, >, >=, <, <=", s)
	}

	op, column, err := parseQueryTerm(parts[0], false)
	if err != nil {
		return api.Having{}, fmt.Errorf("invalid having %q: %w", s, err)
	}
	if op == "HEATMAP" {
		return api.Having{}, fmt.Errorf("invalid having %q: HEATMAP is not supported in havings", s)
	}

	value, err := parseNumericValue(parts[2])
	if err != nil {
		return api.Having{}, fmt.Errorf("invalid having %q: value must be a number", s)
	}

	return api.Having{
		CalculateOp: op,
		Column:      column,
		Op:          parts[1],
		Value:       value,
	}, nil
}

func parseQueryTerm(s string, allowPlainColumn bool) (op string, column string, err error) {
	term := strings.TrimSpace(s)
	if term == "" {
		return "", "", fmt.Errorf("term is empty")
	}

	if open := strings.Index(term, "("); open >= 0 {
		if !strings.HasSuffix(term, ")") || strings.Count(term, "(") != 1 || strings.Count(term, ")") != 1 {
			return "", "", fmt.Errorf("term must use CALC(column) syntax")
		}
		rawOp := strings.ToUpper(strings.TrimSpace(term[:open]))
		if !calculationOps[rawOp] {
			return "", "", fmt.Errorf("unknown calculation op %q", rawOp)
		}
		rawColumn := strings.TrimSpace(term[open+1 : len(term)-1])
		if rawColumn == "" && rawOp != "COUNT" && rawOp != "CONCURRENCY" {
			return "", "", fmt.Errorf("calculation op %s requires a column", rawOp)
		}
		return rawOp, rawColumn, nil
	}

	rawOp := strings.ToUpper(term)
	if calculationOps[rawOp] && rawOp == term {
		return rawOp, "", nil
	}
	if allowPlainColumn {
		return "", term, nil
	}
	return "", "", fmt.Errorf("term must be a calculation like MAX(duration_ms) or COUNT")
}

func parseNumericValue(s string) (any, error) {
	if strings.ContainsAny(s, ".eE") {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func validateQueryReferences(query *api.Query) error {
	breakdowns := make(map[string]bool, len(query.Breakdowns))
	for _, breakdown := range query.Breakdowns {
		breakdowns[breakdown] = true
	}

	calculations := make(map[string]bool, len(query.Calculations))
	for _, calc := range query.Calculations {
		calculations[calculationKey(calc.Op, calc.Column)] = true
	}

	for _, order := range query.Orders {
		if order.Op != "" {
			if !calculations[calculationKey(order.Op, order.Column)] {
				return fmt.Errorf("invalid order %q: calculation must match one of the query calculations", formatOrderTerm(order))
			}
			continue
		}
		if !breakdowns[order.Column] {
			return fmt.Errorf("invalid order %q: column must match one of the query breakdowns", order.Column)
		}
	}

	for _, having := range query.Havings {
		if !calculations[calculationKey(having.CalculateOp, having.Column)] {
			return fmt.Errorf("invalid having %q: calculation must match one of the query calculations", formatHavingTerm(having))
		}
	}

	return nil
}

func calculationKey(op, column string) string {
	return strings.ToUpper(op) + "\x00" + column
}

func formatOrderTerm(order api.Order) string {
	if order.Op == "" {
		return order.Column
	}
	return formatCalculationTerm(order.Op, order.Column)
}

func formatHavingTerm(having api.Having) string {
	return formatCalculationTerm(having.CalculateOp, having.Column)
}

func formatCalculationTerm(op, column string) string {
	if column == "" {
		return strings.ToUpper(op)
	}
	return fmt.Sprintf("%s(%s)", strings.ToUpper(op), column)
}

func GetQueryCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-query",
		Category: "Queries",
		Usage:    "Get a query by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Query ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			query, err := client.GetQuery(ctx, cmd.String("dataset"), cmd.String("id"))
			if err != nil {
				return err
			}

			return printJSON(query)
		},
	}
}

func queryBuildFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "dataset",
			Usage:    "Dataset slug",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "calculation-op",
			Usage: "Calculation operation (e.g. COUNT, AVG, P99); repeat for multiple calculations",
		},
		&cli.StringSliceFlag{
			Name:  "calculation-column",
			Usage: "Calculation column (use empty string for COUNT); repeat to match each --calculation-op",
		},
		&cli.StringSliceFlag{
			Name:  "breakdown",
			Usage: "Breakdown column; repeat for multiple dimensions",
		},
		&cli.StringSliceFlag{
			Name:  "filter",
			Usage: `Filter in "column op [value]" form; repeat for multiple filters (e.g. --filter "duration_ms > 100" --filter "name exists")`,
		},
		&cli.StringFlag{
			Name:  "filter-combination",
			Usage: "How to combine filters: AND (default) or OR",
		},
		&cli.StringSliceFlag{
			Name:  "order",
			Usage: `Order in "term asc|desc" form; repeat for multiple orders (e.g. --order "MAX(duration_ms) desc" --order "trace.trace_id asc")`,
		},
		&cli.IntFlag{
			Name:  "limit",
			Usage: "Maximum number of unique groups returned in results (1-1000)",
		},
		&cli.StringSliceFlag{
			Name:  "having",
			Usage: `Having clause in "CALC(column) op number" form; repeat for multiple havings (e.g. --having "MAX(duration_ms) > 1000")`,
		},
		&cli.StringFlag{
			Name:  "time-range",
			Usage: `Time range (e.g. 3600, "4 hours", "last week")`,
		},
		&cli.StringFlag{
			Name:  "from",
			Usage: `Start time (e.g. "2024-02-11 18:00", "2024-02-11T18:00:00Z")`,
		},
		&cli.StringFlag{
			Name:  "to",
			Usage: `End time (e.g. "2024-02-11 18:45", "2024-02-11T18:00:00Z")`,
		},
		&cli.StringFlag{
			Name:  "timezone",
			Usage: `Timezone for parsing dates (e.g. "America/New_York", default UTC)`,
		},
		&cli.StringFlag{
			Name:  "query-json",
			Usage: "Path to raw Honeycomb query JSON, or - for stdin; cannot be combined with query-building flags",
		},
	}
}

func queryBuildDescription() string {
	return `Examples:
  hccli create-query --dataset aws --calculation-op COUNT --breakdown service.name --order "COUNT desc" --limit 10
  hccli create-query --dataset aws --calculation-op MAX --calculation-column duration_ms --breakdown trace.trace_id --order "MAX(duration_ms) desc" --limit 1
  hccli create-query --dataset aws --calculation-op MAX --calculation-column duration_ms --having "MAX(duration_ms) > 1000"`
}

func buildQueryFromFlags(cmd *cli.Command) (*api.Query, error) {
	ops := cmd.StringSlice("calculation-op")
	cols := cmd.StringSlice("calculation-column")

	if len(ops) == 0 {
		return nil, fmt.Errorf("one or more --calculation-op values are required unless --query-json is used")
	}

	if len(cols) > 0 && len(cols) != len(ops) {
		return nil, fmt.Errorf("number of --calculation-column values (%d) must match --calculation-op values (%d)", len(cols), len(ops))
	}

	var calcs []api.Calculation
	for i, op := range ops {
		c := api.Calculation{Op: op}
		if i < len(cols) && cols[i] != "" {
			c.Column = cols[i]
		}
		calcs = append(calcs, c)
	}

	query := &api.Query{
		Calculations: calcs,
	}

	if v := cmd.StringSlice("breakdown"); len(v) > 0 {
		query.Breakdowns = v
	}

	for _, raw := range cmd.StringSlice("filter") {
		f, err := parseFilter(raw)
		if err != nil {
			return nil, err
		}
		query.Filters = append(query.Filters, f)
	}

	if v := cmd.String("filter-combination"); v != "" {
		query.FilterCombination = v
	}

	for _, raw := range cmd.StringSlice("order") {
		order, err := parseOrder(raw)
		if err != nil {
			return nil, err
		}
		query.Orders = append(query.Orders, order)
	}

	if cmd.IsSet("limit") {
		limit := cmd.Int("limit")
		if limit < 1 || limit > 1000 {
			return nil, fmt.Errorf("invalid limit %d: must be between 1 and 1000", limit)
		}
		query.Limit = limit
	}

	for _, raw := range cmd.StringSlice("having") {
		having, err := parseHaving(raw)
		if err != nil {
			return nil, err
		}
		query.Havings = append(query.Havings, having)
	}

	if err := validateQueryReferences(query); err != nil {
		return nil, err
	}

	loc := time.UTC
	if tz := cmd.String("timezone"); tz != "" {
		var err error
		loc, err = time.LoadLocation(tz)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %q: %w", tz, err)
		}
	}

	if v := cmd.String("time-range"); v != "" {
		tr, err := timefmt.ParseTimeRange(v)
		if err != nil {
			return nil, fmt.Errorf("invalid time-range %q: %w", v, err)
		}
		query.TimeRange = tr
	}

	if v := cmd.String("from"); v != "" {
		ts, err := timefmt.ParseTimestamp(v, loc)
		if err != nil {
			return nil, fmt.Errorf("invalid from time %q: %w", v, err)
		}
		query.StartTime = int(ts)
	}

	if v := cmd.String("to"); v != "" {
		ts, err := timefmt.ParseTimestamp(v, loc)
		if err != nil {
			return nil, fmt.Errorf("invalid to time %q: %w", v, err)
		}
		query.EndTime = int(ts)
	}

	return query, nil
}

func buildQueryInputFromFlags(cmd *cli.Command) (*queryInput, error) {
	if path := cmd.String("query-json"); path != "" {
		for _, name := range queryFlagNames {
			if cmd.IsSet(name) {
				return nil, fmt.Errorf("--query-json cannot be combined with --%s", name)
			}
		}

		raw, err := readQueryJSON(path)
		if err != nil {
			return nil, err
		}
		return &queryInput{RawJSON: raw}, nil
	}

	query, err := buildQueryFromFlags(cmd)
	if err != nil {
		return nil, err
	}
	return &queryInput{Query: query}, nil
}

func readQueryJSON(path string) ([]byte, error) {
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("reading query-json from stdin: %w", err)
		}
	} else {
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading query-json %q: %w", path, err)
		}
	}

	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("invalid query JSON %q: %w", path, err)
	}
	return bytes.TrimSpace(data), nil
}

func CreateQueryCmd() *cli.Command {
	return &cli.Command{
		Name:        "create-query",
		Category:    "Queries",
		Usage:       "Create a new query",
		Description: queryBuildDescription(),
		Flags:       queryBuildFlags(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			input, err := buildQueryInputFromFlags(cmd)
			if err != nil {
				return err
			}

			if input.RawJSON != nil {
				created, err := client.CreateQueryRaw(ctx, cmd.String("dataset"), input.RawJSON)
				if err != nil {
					return err
				}
				return printJSON(created)
			}

			created, err := client.CreateQuery(ctx, cmd.String("dataset"), input.Query)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}

func RunQueryCmd() *cli.Command {
	flags := queryBuildFlags()
	flags = append(flags, queryResultPollingFlags()...)

	return &cli.Command{
		Name:     "run-query",
		Category: "Query Results",
		Usage:    "Create a query and return results (polls until complete)",
		Description: `Examples:
  hccli run-query --dataset aws --calculation-op COUNT --breakdown service.name --order "COUNT desc" --limit 10
  hccli run-query --dataset aws --calculation-op MAX --calculation-column duration_ms --filter "http.route contains /service/awards" --time-range "30 minutes"`,
		Flags: flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			input, err := buildQueryInputFromFlags(cmd)
			if err != nil {
				return err
			}

			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			dataset := cmd.String("dataset")
			var queryID string
			if input.RawJSON != nil {
				created, err := client.CreateQueryRaw(ctx, dataset, input.RawJSON)
				if err != nil {
					return err
				}
				queryID, _ = created["id"].(string)
			} else {
				created, err := client.CreateQuery(ctx, dataset, input.Query)
				if err != nil {
					return err
				}
				queryID = created.ID
			}
			if queryID == "" {
				return fmt.Errorf("created query response did not include an id")
			}

			pollInterval := time.Duration(cmd.Int("poll-interval")) * time.Second
			timeout := time.Duration(cmd.Int("timeout")) * time.Second
			result, err := pollQueryResult(ctx, client, dataset, queryID, pollInterval, timeout)
			if err != nil {
				return err
			}

			return printJSON(result)
		},
	}
}
