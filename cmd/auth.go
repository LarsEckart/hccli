package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

type authListOutput struct {
	ActiveProfile string           `json:"active_profile,omitzero"`
	LocalProfile  string           `json:"local_profile,omitzero"`
	ConfigPath    string           `json:"config_path"`
	LocalPath     string           `json:"local_path"`
	Profiles      []profileSummary `json:"profiles"`
}

type profileSummary struct {
	Name            string `json:"name"`
	Active          bool   `json:"active"`
	LocalActive     bool   `json:"local_active"`
	APIURL          string `json:"api_url,omitzero"`
	KeyID           string `json:"key_id,omitzero"`
	KeyType         string `json:"key_type,omitzero"`
	TeamName        string `json:"team_name,omitzero"`
	TeamSlug        string `json:"team_slug,omitzero"`
	EnvironmentName string `json:"environment_name,omitzero"`
	EnvironmentSlug string `json:"environment_slug,omitzero"`
	UpdatedAt       string `json:"updated_at,omitzero"`
}

type authStatusOutput struct {
	Source  string            `json:"source"`
	Profile string            `json:"profile,omitzero"`
	APIURL  string            `json:"api_url"`
	Valid   *bool             `json:"valid,omitzero"`
	Info    *profileSummary   `json:"profile_info,omitzero"`
	Auth    *api.AuthResponse `json:"auth,omitzero"`
}

type authLoginOutput struct {
	Profile     string         `json:"profile"`
	Active      bool           `json:"active"`
	LocalActive bool           `json:"local_active"`
	ConfigPath  string         `json:"config_path"`
	LocalPath   string         `json:"local_path,omitzero"`
	Info        profileSummary `json:"profile_info"`
}

type authSwitchOutput struct {
	Profile     string `json:"profile"`
	Active      bool   `json:"active"`
	LocalActive bool   `json:"local_active"`
	ConfigPath  string `json:"config_path,omitzero"`
	LocalPath   string `json:"local_path,omitzero"`
}

type authLogoutOutput struct {
	Profile       string `json:"profile"`
	Removed       bool   `json:"removed"`
	ActiveCleared bool   `json:"active_cleared"`
	LocalCleared  bool   `json:"local_cleared"`
	ConfigPath    string `json:"config_path"`
}

func AuthCmd() *cli.Command {
	return &cli.Command{
		Name:     "auth",
		Category: "Auth",
		Usage:    "Manage Honeycomb authentication profiles",
		Description: `Manage named Honeycomb profiles, similar to gh auth switch.

Run hccli auth status to see which credentials will be used.
Run hccli auth login --profile <name> --api-key-stdin to add a profile.`,
		Action: authWhoamiAction,
		Commands: []*cli.Command{
			authStatusCmd(),
			authWhoamiCmd(),
			authWhoamiV2Cmd(),
			authLoginCmd(),
			authListCmd(),
			authSwitchCmd(),
			authLogoutCmd(),
		},
	}
}

func AuthV2Cmd() *cli.Command {
	return &cli.Command{
		Name:     "auth-v2",
		Category: "Auth",
		Usage:    "Show management API key info and permissions (v2)",
		Action:   authWhoamiV2Action,
	}
}

func authStatusCmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show the active Honeycomb credentials and validate them",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-verify",
				Usage: "Do not call Honeycomb to validate credentials",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			credentials, err := resolveCredentials(cmd)
			if err != nil {
				return err
			}

			out := authStatusOutput{
				Source:  credentials.Source,
				Profile: credentials.Profile,
				APIURL:  credentials.APIURL,
			}

			if credentials.Profile != "" {
				cfg, err := loadConfig()
				if err != nil {
					return err
				}
				if profile, ok := cfg.Profiles[credentials.Profile]; ok {
					summary := summarizeProfile(credentials.Profile, profile, cfg.ActiveProfile, localActiveProfile())
					out.Info = &summary
				}
			}

			if !cmd.Bool("no-verify") {
				client, err := newClient(cmd)
				if err != nil {
					return err
				}
				auth, err := client.GetAuth(ctx)
				valid := err == nil
				out.Valid = &valid
				if err != nil {
					return fmt.Errorf("validating credentials: %w", err)
				}
				out.Auth = auth
			}

			return printJSON(out)
		},
	}
}

func authWhoamiCmd() *cli.Command {
	return &cli.Command{
		Name:   "whoami",
		Usage:  "Show API key info and permissions",
		Action: authWhoamiAction,
	}
}

func authWhoamiV2Cmd() *cli.Command {
	return &cli.Command{
		Name:   "whoami-v2",
		Usage:  "Show management API key info and permissions (v2)",
		Action: authWhoamiV2Action,
	}
}

func authLoginCmd() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Store a Honeycomb API key as a named profile",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "Profile name to create or update",
			},
			&cli.StringFlag{
				Name:    "api-key",
				Sources: cli.EnvVars("HONEYCOMB_API_KEY"),
				Usage:   "Honeycomb API key; prefer --api-key-stdin to avoid shell history",
			},
			&cli.BoolFlag{
				Name:  "api-key-stdin",
				Usage: "Read the Honeycomb API key from stdin",
			},
			&cli.BoolFlag{
				Name:  "skip-verify",
				Usage: "Store the key without validating it against Honeycomb",
			},
			&cli.BoolFlag{
				Name:  "no-switch",
				Usage: "Do not make this profile active after login",
			},
			&cli.BoolFlag{
				Name:  "local",
				Usage: "Make this profile active only for the current project",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			name := strings.TrimSpace(cmd.String("profile"))
			if err := validateProfileName(name); err != nil {
				return err
			}
			if cmd.Bool("local") && cmd.Bool("no-switch") {
				return fmt.Errorf("--local and --no-switch cannot be used together")
			}

			apiKey, err := loginAPIKey(cmd)
			if err != nil {
				return err
			}

			profile := Profile{
				APIKey:    apiKey,
				APIURL:    cmd.String("api-url"),
				UpdatedAt: nowTimestamp(),
			}

			if !cmd.Bool("skip-verify") {
				client := api.NewClient(apiKey, time.Duration(cmd.Int("timeout"))*time.Second)
				client.BaseURL = cmd.String("api-url")
				auth, err := client.GetAuth(ctx)
				if err != nil {
					return fmt.Errorf("validating API key: %w", err)
				}
				profile.TeamName = auth.Team.Name
				profile.TeamSlug = auth.Team.Slug
				profile.EnvironmentName = auth.Environment.Name
				profile.EnvironmentSlug = auth.Environment.Slug
				profile.KeyID = auth.ID
				profile.KeyType = auth.Type
			}

			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			cfg.Profiles[name] = profile
			if !cmd.Bool("no-switch") && !cmd.Bool("local") {
				cfg.ActiveProfile = name
			}
			if err := saveConfig(cfg); err != nil {
				return err
			}

			var localPath string
			if cmd.Bool("local") {
				local := &LocalConfig{ActiveProfile: name}
				if err := saveLocalConfig(local); err != nil {
					return err
				}
				localPath, err = localConfigPath()
				if err != nil {
					return err
				}
			}

			configPath, err := configPath()
			if err != nil {
				return err
			}

			return printJSON(authLoginOutput{
				Profile:     name,
				Active:      cfg.ActiveProfile == name,
				LocalActive: cmd.Bool("local"),
				ConfigPath:  configPath,
				LocalPath:   localPath,
				Info:        summarizeProfile(name, profile, cfg.ActiveProfile, localActiveProfile()),
			})
		},
	}
}

func authListCmd() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List configured Honeycomb profiles",
		Action: func(context.Context, *cli.Command) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			localProfile := localActiveProfile()
			profiles := make([]profileSummary, 0, len(cfg.Profiles))
			for _, name := range profileNames(cfg) {
				profiles = append(profiles, summarizeProfile(name, cfg.Profiles[name], cfg.ActiveProfile, localProfile))
			}
			configPath, err := configPath()
			if err != nil {
				return err
			}
			localPath, err := localConfigPath()
			if err != nil {
				return err
			}
			return printJSON(authListOutput{
				ActiveProfile: cfg.ActiveProfile,
				LocalProfile:  localProfile,
				ConfigPath:    configPath,
				LocalPath:     localPath,
				Profiles:      profiles,
			})
		},
	}
}

func authSwitchCmd() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Usage:     "Switch the active Honeycomb profile",
		ArgsUsage: "<profile>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "local",
				Usage: "Switch only for the current project",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			name := strings.TrimSpace(cmd.Args().First())
			if err := validateProfileName(name); err != nil {
				return err
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if _, ok := cfg.Profiles[name]; !ok {
				return fmt.Errorf("honeycomb profile %q not found; run `hccli auth list` to see configured profiles", name)
			}

			out := authSwitchOutput{Profile: name}
			if cmd.Bool("local") {
				if err := saveLocalConfig(&LocalConfig{ActiveProfile: name}); err != nil {
					return err
				}
				localPath, err := localConfigPath()
				if err != nil {
					return err
				}
				out.LocalActive = true
				out.LocalPath = localPath
			} else {
				cfg.ActiveProfile = name
				if err := saveConfig(cfg); err != nil {
					return err
				}
				configPath, err := configPath()
				if err != nil {
					return err
				}
				out.Active = true
				out.ConfigPath = configPath
			}
			return printJSON(out)
		},
	}
}

func authLogoutCmd() *cli.Command {
	return &cli.Command{
		Name:      "logout",
		Usage:     "Remove a stored Honeycomb profile",
		ArgsUsage: "<profile>",
		Action: func(_ context.Context, cmd *cli.Command) error {
			name := strings.TrimSpace(cmd.Args().First())
			if err := validateProfileName(name); err != nil {
				return err
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			_, existed := cfg.Profiles[name]
			delete(cfg.Profiles, name)
			activeCleared := false
			if cfg.ActiveProfile == name {
				cfg.ActiveProfile = ""
				activeCleared = true
			}
			if err := saveConfig(cfg); err != nil {
				return err
			}
			localCleared := false
			local, err := loadLocalConfig()
			if err != nil {
				return err
			}
			if local.ActiveProfile == name {
				local.ActiveProfile = ""
				if err := saveLocalConfig(local); err != nil {
					return err
				}
				localCleared = true
			}
			configPath, err := configPath()
			if err != nil {
				return err
			}
			return printJSON(authLogoutOutput{
				Profile:       name,
				Removed:       existed,
				ActiveCleared: activeCleared,
				LocalCleared:  localCleared,
				ConfigPath:    configPath,
			})
		},
	}
}

func authWhoamiAction(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	auth, err := client.GetAuth(ctx)
	if err != nil {
		return err
	}

	return printJSON(auth)
}

func authWhoamiV2Action(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	auth, err := client.GetAuthV2(ctx)
	if err != nil {
		return err
	}

	return printJSON(auth)
}

func loginAPIKey(cmd *cli.Command) (string, error) {
	if cmd.Bool("api-key-stdin") {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading API key from stdin: %w", err)
		}
		apiKey := strings.TrimSpace(string(data))
		if apiKey == "" {
			return "", fmt.Errorf("no API key found on stdin")
		}
		return apiKey, nil
	}

	apiKey := strings.TrimSpace(cmd.String("api-key"))
	if apiKey == "" {
		return "", fmt.Errorf("api-key is required; pass --api-key-stdin, --api-key, or set HONEYCOMB_API_KEY")
	}
	return apiKey, nil
}

func validateProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("profile is required")
	}
	if strings.ContainsAny(name, `/\\`) || strings.Contains(name, "..") || strings.ContainsAny(name, " \t\n\r") {
		return fmt.Errorf("profile %q is invalid; use a simple name like work or personal", name)
	}
	return nil
}

func summarizeProfile(name string, profile Profile, activeProfile string, localProfile string) profileSummary {
	return profileSummary{
		Name:            name,
		Active:          activeProfile == name,
		LocalActive:     localProfile == name,
		APIURL:          profile.APIURL,
		KeyID:           profile.KeyID,
		KeyType:         profile.KeyType,
		TeamName:        profile.TeamName,
		TeamSlug:        profile.TeamSlug,
		EnvironmentName: profile.EnvironmentName,
		EnvironmentSlug: profile.EnvironmentSlug,
		UpdatedAt:       profile.UpdatedAt,
	}
}

func localActiveProfile() string {
	local, err := loadLocalConfig()
	if err != nil {
		return ""
	}
	return local.ActiveProfile
}
