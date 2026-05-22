package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

const commandHelpWithoutGlobalOptionsTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{template "usageTemplate" .}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`

var helpExamples = map[string]string{
	"hccli": `Examples:
  hccli auth status
  hccli datasets
  hccli run-query --dataset prod --calculation-op COUNT --time-range "30 minutes"`,
	"hccli help": `Examples:
  hccli help
  hccli help run-query
  hccli help auth login`,
	"hccli auth": `Examples:
  hccli auth status
  hccli auth login --profile work --api-key-stdin
  hccli auth switch personal`,
	"hccli auth status": `Examples:
  hccli auth status
  hccli auth status --no-verify`,
	"hccli auth whoami": `Examples:
  hccli auth whoami
  hccli --profile work auth whoami`,
	"hccli auth whoami-v2": `Examples:
  hccli auth whoami-v2
  hccli --profile work auth whoami-v2`,
	"hccli auth login": `Global --api-url and --timeout may be passed before auth login when validating the key.

Examples:
  printf '%s\n' "$HONEYCOMB_API_KEY" | hccli auth login --profile work --api-key-stdin
  hccli --api-url https://api.honeycomb.io auth login --profile work --api-key-stdin --local
  hccli auth login --profile ci --api-key-stdin --no-switch`,
	"hccli auth list": `Examples:
  hccli auth list
  hccli auth ls`,
	"hccli auth switch": `Examples:
  hccli auth switch work
  hccli auth switch personal --local`,
	"hccli auth logout": `Examples:
  hccli auth logout old-profile`,
	"hccli auth-v2": `Alias for hccli auth whoami-v2.

Examples:
  hccli auth-v2
  hccli --profile work auth-v2`,
	"hccli boards": `Examples:
  hccli boards`,
	"hccli get-board": `Examples:
  hccli get-board --id brd_123`,
	"hccli create-board": `Examples:
  hccli create-board --name "Production Overview" --description "Service health dashboard"`,
	"hccli update-board": `Examples:
  hccli update-board --id brd_123 --name "Production Overview" --description "Updated description"
  hccli update-board --id brd_123 --query-id qry_123 --query-caption "Latency by route"`,
	"hccli delete-board": `Examples:
  hccli delete-board --id brd_123`,
	"hccli board-views": `Examples:
  hccli board-views --board-id brd_123`,
	"hccli get-board-view": `Examples:
  hccli get-board-view --board-id brd_123 --view-id view_123`,
	"hccli create-board-view": `Examples:
  hccli create-board-view --board-id brd_123 --name "Errors" --query-id qry_123 --query-caption "Errors by service"`,
	"hccli update-board-view": `Examples:
  hccli update-board-view --board-id brd_123 --view-id view_123 --name "Errors" --query-id qry_123 --query-caption "Errors by service"`,
	"hccli delete-board-view": `Examples:
  hccli delete-board-view --board-id brd_123 --view-id view_123`,
	"hccli burn-alerts": `Examples:
  hccli burn-alerts --dataset prod --slo-id slo_123`,
	"hccli get-burn-alert": `Examples:
  hccli get-burn-alert --dataset prod --id ba_123`,
	"hccli update-burn-alert": `Examples:
  hccli update-burn-alert --dataset prod --id ba_123 --recipient type=email,target=alerts@example.com --exhaustion-minutes 60`,
	"hccli delete-burn-alert": `Examples:
  hccli delete-burn-alert --dataset prod --id ba_123`,
	"hccli columns": `Examples:
  hccli columns --dataset prod`,
	"hccli get-column": `Examples:
  hccli get-column --dataset prod --id col_123`,
	"hccli create-column": `Examples:
  hccli create-column --dataset prod --key_name service.name --type string --description "Service name"`,
	"hccli update-column": `Examples:
  hccli update-column --dataset prod --id col_123 --description "Service name"`,
	"hccli delete-column": `Examples:
  hccli delete-column --dataset prod --id col_123`,
	"hccli datasets": `Examples:
  hccli datasets`,
	"hccli get-dataset": `Examples:
  hccli get-dataset --slug prod`,
	"hccli create-dataset": `Examples:
  hccli create-dataset --name prod --description "Production telemetry"`,
	"hccli update-dataset": `Examples:
  hccli update-dataset --slug prod --description "Production telemetry" --delete-protected=true`,
	"hccli delete-dataset": `Examples:
  hccli delete-dataset --slug old-dataset`,
	"hccli derived-columns": `Examples:
  hccli derived-columns --dataset prod`,
	"hccli get-derived-column": `Examples:
  hccli get-derived-column --dataset prod --id dc_123`,
	"hccli create-derived-column": `Examples:
  hccli create-derived-column --dataset prod --alias is_error --expression "BOOL(500 <= $status_code)"`,
	"hccli update-derived-column": `Examples:
  hccli update-derived-column --dataset prod --id dc_123 --alias is_error --expression "BOOL(500 <= $status_code)"`,
	"hccli delete-derived-column": `Examples:
  hccli delete-derived-column --dataset prod --id dc_123`,
	"hccli marker-settings": `Examples:
  hccli marker-settings --dataset prod`,
	"hccli create-marker-setting": `Examples:
  hccli create-marker-setting --dataset prod --type deploy --color "#1f77b4"`,
	"hccli update-marker-setting": `Examples:
  hccli update-marker-setting --dataset prod --id ms_123 --type deploy --color "#2ca02c"`,
	"hccli delete-marker-setting": `Examples:
  hccli delete-marker-setting --dataset prod --id ms_123`,
	"hccli markers": `Examples:
  hccli markers --dataset prod`,
	"hccli create-marker": `Examples:
  hccli create-marker --dataset prod --message "Deploy v1.2.3" --type deploy --url https://github.com/example/repo/actions/runs/123`,
	"hccli update-marker": `Examples:
  hccli update-marker --dataset prod --id mrk_123 --message "Deploy v1.2.4" --type deploy`,
	"hccli delete-marker": `Examples:
  hccli delete-marker --dataset prod --id mrk_123`,
	"hccli get-query": `Examples:
  hccli get-query --dataset prod --id qry_123`,
	"hccli create-query-result": `Examples:
  hccli create-query-result --dataset prod --query-id qry_123 --poll-interval 2 --result-timeout 60`,
	"hccli get-query-result": `Examples:
  hccli get-query-result --dataset prod --id qres_123`,
	"hccli query-annotations": `Examples:
  hccli query-annotations --dataset prod`,
	"hccli create-query-annotation": `Examples:
  hccli create-query-annotation --dataset prod --name "Latency by route" --query-id qry_123`,
	"hccli get-query-annotation": `Examples:
  hccli get-query-annotation --dataset prod --id qa_123`,
	"hccli update-query-annotation": `Examples:
  hccli update-query-annotation --dataset prod --id qa_123 --name "Latency by route" --query-id qry_123`,
	"hccli delete-query-annotation": `Examples:
  hccli delete-query-annotation --dataset prod --id qa_123`,
	"hccli slos": `Examples:
  hccli slos --dataset prod`,
	"hccli get-slo": `Examples:
  hccli get-slo --dataset prod --id slo_123`,
	"hccli create-slo": `Examples:
  hccli create-slo --dataset prod --name "API availability" --sli-alias api_available --time-period-days 30 --target-per-million 999000`,
	"hccli update-slo": `Examples:
  hccli update-slo --dataset prod --id slo_123 --name "API availability" --sli-alias api_available --time-period-days 30 --target-per-million 999000`,
	"hccli delete-slo": `Examples:
  hccli delete-slo --dataset prod --id slo_123`,
	"hccli get-trace": `Examples:
  hccli get-trace --dataset prod --trace-id 2f1c8f0b9a123456`,
}

// ApplyHelpConventions adds machine-friendly help details that urfave/cli does
// not render by default: required markers, examples, and hidden zero defaults
// for optional integer flags where zero only means "unset".
func ApplyHelpConventions(root *cli.Command) {
	applyHelpConventions(root, root.Name)
}

func applyHelpConventions(command *cli.Command, fullName string) {
	for _, flag := range command.Flags {
		applyFlagHelpConventions(flag)
	}

	if command.Name == "login" && strings.HasSuffix(fullName, " auth login") {
		command.CustomHelpTemplate = commandHelpWithoutGlobalOptionsTemplate
	}

	if examples, ok := helpExamples[fullName]; ok && !strings.Contains(command.Description, "Examples:") {
		command.Description = appendDescription(command.Description, examples)
	}
	if strings.HasPrefix(command.Name, "delete-") && !strings.Contains(command.Description, "Deletes immediately") {
		command.Description = appendDescription("Deletes immediately; no confirmation prompt.", command.Description)
	}

	for _, subcommand := range command.Commands {
		applyHelpConventions(subcommand, fullName+" "+subcommand.Name)
	}
}

func applyFlagHelpConventions(flag cli.Flag) {
	if required, ok := flag.(cli.RequiredFlag); ok && required.IsRequired() {
		appendRequiredToUsage(flag)
	}

	if intFlag, ok := flag.(*cli.IntFlag); ok && !intFlag.Required && intFlag.Value == 0 && intFlag.DefaultText == "" {
		intFlag.HideDefault = true
	}
}

func appendRequiredToUsage(flag cli.Flag) {
	switch f := flag.(type) {
	case *cli.StringFlag:
		f.Usage = appendUsageSuffix(f.Usage, "(required)")
	case *cli.IntFlag:
		f.Usage = appendUsageSuffix(f.Usage, "(required)")
	case *cli.Int64Flag:
		f.Usage = appendUsageSuffix(f.Usage, "(required)")
	case *cli.BoolFlag:
		f.Usage = appendUsageSuffix(f.Usage, "(required)")
	case *cli.StringSliceFlag:
		f.Usage = appendUsageSuffix(f.Usage, "(required)")
	}
}

func appendUsageSuffix(usage string, suffix string) string {
	if strings.Contains(usage, suffix) {
		return usage
	}
	if usage == "" {
		return suffix
	}
	return usage + " " + suffix
}

func appendDescription(description string, addition string) string {
	description = strings.TrimSpace(description)
	addition = strings.TrimSpace(addition)
	if description == "" {
		return addition
	}
	if addition == "" {
		return description
	}
	return description + "\n\n" + addition
}

func HelpCmd() *cli.Command {
	return &cli.Command{
		Name:      "help",
		Aliases:   []string{"h"},
		Usage:     "Shows a list of commands or help for one command",
		ArgsUsage: "[command]",
		Action: func(_ context.Context, command *cli.Command) error {
			args := command.Args().Slice()
			if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
				return cli.ShowRootCommandHelp(command.Root())
			}

			current := command.Root()
			for _, arg := range args {
				next := findSubcommand(current, arg)
				if next == nil {
					return fmt.Errorf("no help topic for %q", strings.Join(args, " "))
				}
				current = next
			}

			if len(current.VisibleCommands()) > 0 {
				return cli.ShowSubcommandHelp(current)
			}
			template := current.CustomHelpTemplate
			if template == "" {
				template = cli.CommandHelpTemplate
			}
			cli.HelpPrinter(current.Root().Writer, template, current)
			return nil
		},
	}
}

func findSubcommand(command *cli.Command, name string) *cli.Command {
	for _, subcommand := range command.Commands {
		if subcommand.HasName(name) {
			return subcommand
		}
	}
	return nil
}
