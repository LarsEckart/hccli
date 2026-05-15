package cmd

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

const (
	defaultAPIURL      = "https://api.honeycomb.io"
	configDirEnv       = "HCCLI_CONFIG_DIR"
	configFileName     = "config.json"
	localConfigDirName = ".hccli"
)

type Config struct {
	ActiveProfile string             `json:"active_profile,omitzero"`
	Profiles      map[string]Profile `json:"profiles,omitzero"`
}

type LocalConfig struct {
	ActiveProfile string `json:"active_profile,omitzero"`
}

type Profile struct {
	APIKey          string `json:"api_key,omitzero"`
	APIURL          string `json:"api_url,omitzero"`
	KeyID           string `json:"key_id,omitzero"`
	KeyType         string `json:"key_type,omitzero"`
	TeamName        string `json:"team_name,omitzero"`
	TeamSlug        string `json:"team_slug,omitzero"`
	EnvironmentName string `json:"environment_name,omitzero"`
	EnvironmentSlug string `json:"environment_slug,omitzero"`
	UpdatedAt       string `json:"updated_at,omitzero"`
}

type resolvedCredentials struct {
	APIKey  string
	APIURL  string
	Profile string
	Source  string
}

func resolveCredentials(cmd *cli.Command) (*resolvedCredentials, error) {
	apiKey := strings.TrimSpace(cmd.String("api-key"))
	if apiKey != "" {
		return &resolvedCredentials{
			APIKey: apiKey,
			APIURL: cmd.String("api-url"),
			Source: "api-key",
		}, nil
	}

	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	profileName := strings.TrimSpace(cmd.String("profile"))
	source := "profile"
	if profileName == "" {
		local, err := loadLocalConfig()
		if err != nil {
			return nil, err
		}
		if local.ActiveProfile != "" {
			profileName = local.ActiveProfile
			source = "local-profile"
		}
	}
	if profileName == "" && cfg.ActiveProfile != "" {
		profileName = cfg.ActiveProfile
		source = "active-profile"
	}
	if profileName == "" {
		return nil, errors.New("no Honeycomb API key configured; run `hccli auth login --profile <name> --api-key-stdin`, set HONEYCOMB_API_KEY, or pass --api-key")
	}

	profile, ok := cfg.Profiles[profileName]
	if !ok {
		return nil, fmt.Errorf("honeycomb profile %q not found; run `hccli auth list` to see configured profiles", profileName)
	}
	if strings.TrimSpace(profile.APIKey) == "" {
		return nil, fmt.Errorf("honeycomb profile %q has no API key; run `hccli auth login --profile %s --api-key-stdin` to repair it", profileName, profileName)
	}

	return &resolvedCredentials{
		APIKey:  profile.APIKey,
		APIURL:  configuredAPIURL(cmd, profile),
		Profile: profileName,
		Source:  source,
	}, nil
}

func configuredAPIURL(cmd *cli.Command, profile Profile) string {
	if cmd.IsSet("api-url") {
		return cmd.String("api-url")
	}
	return cmp.Or(profile.APIURL, cmd.String("api-url"), defaultAPIURL)
}

func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return newConfig(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("securing config permissions: %w", err)
	}
	return nil
}

func newConfig() *Config {
	return &Config{Profiles: make(map[string]Profile)}
}

func configPath() (string, error) {
	dir, err := userConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

func userConfigDir() (string, error) {
	if dir := strings.TrimSpace(os.Getenv(configDirEnv)); dir != "" {
		return dir, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("finding user config directory: %w", err)
	}
	return filepath.Join(dir, "hccli"), nil
}

func loadLocalConfig() (*LocalConfig, error) {
	path, err := localConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &LocalConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading local config: %w", err)
	}
	var cfg LocalConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing local config %s: %w", path, err)
	}
	return &cfg, nil
}

func saveLocalConfig(cfg *LocalConfig) error {
	path, err := localConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating local config directory: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding local config: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing local config: %w", err)
	}
	return nil
}

func localConfigPath() (string, error) {
	root, err := projectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, localConfigDirName, configFileName), nil
}

func projectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding working directory: %w", err)
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("checking git root: %w", err)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return wd, nil
		}
	}
}

func profileNames(cfg *Config) []string {
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func nowTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
