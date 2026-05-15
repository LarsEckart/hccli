package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeHoneycombAccount struct {
	key      string
	id       string
	keyType  string
	team     string
	env      string
	teamName string
	envName  string
}

func TestAuthProfiles_LoginSwitchAndOverride(t *testing.T) {
	server := newAuthServer(t, []fakeHoneycombAccount{
		{key: "work-key", id: "work-id", keyType: "configuration", team: "work-team", env: "prod", teamName: "Work", envName: "Production"},
		{key: "personal-key", id: "personal-id", keyType: "configuration", team: "personal-team", env: "playground", teamName: "Personal", envName: "Playground"},
	})
	defer server.Close()

	env := isolatedConfigEnv(t)

	stdout, stderr, code := runCLIWithOptions(t, cliRunOptions{Env: env, Stdin: "work-key\n"},
		"--api-url", server.URL,
		"auth", "login", "--profile", "work", "--api-key-stdin",
	)
	if code != 0 {
		t.Fatalf("login work failed with code %d\nstderr: %s", code, stderr)
	}
	login := parseJSON(t, stdout)
	if login["profile"] != "work" || login["active"] != true {
		t.Fatalf("expected work profile to be active, got: %s", stdout)
	}
	assertProfileInfo(t, login["profile_info"], "work-team", "prod")

	_, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env, Stdin: "personal-key\n"},
		"--api-url", server.URL,
		"auth", "login", "--profile", "personal", "--api-key-stdin", "--no-switch",
	)
	if code != 0 {
		t.Fatalf("login personal failed with code %d\nstderr: %s", code, stderr)
	}

	stdout, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env}, "auth", "whoami")
	if code != 0 {
		t.Fatalf("whoami with active profile failed with code %d\nstderr: %s", code, stderr)
	}
	assertAuthTeam(t, stdout, "work-team")

	stdout, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env}, "--profile", "personal", "auth", "whoami")
	if code != 0 {
		t.Fatalf("whoami with --profile failed with code %d\nstderr: %s", code, stderr)
	}
	assertAuthTeam(t, stdout, "personal-team")

	stdout, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env}, "--api-key", "personal-key", "--api-url", server.URL, "auth", "whoami")
	if code != 0 {
		t.Fatalf("whoami with --api-key failed with code %d\nstderr: %s", code, stderr)
	}
	assertAuthTeam(t, stdout, "personal-team")

	stdout, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env}, "auth", "switch", "personal")
	if code != 0 {
		t.Fatalf("switch failed with code %d\nstderr: %s", code, stderr)
	}
	switchOut := parseJSON(t, stdout)
	if switchOut["profile"] != "personal" || switchOut["active"] != true {
		t.Fatalf("expected personal to be active, got: %s", stdout)
	}

	stdout, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env}, "auth", "list")
	if code != 0 {
		t.Fatalf("list failed with code %d\nstderr: %s", code, stderr)
	}
	list := parseJSON(t, stdout)
	if list["active_profile"] != "personal" {
		t.Fatalf("expected personal active profile, got: %s", stdout)
	}
	profiles, ok := list["profiles"].([]any)
	if !ok || len(profiles) != 2 {
		t.Fatalf("expected two profiles, got: %s", stdout)
	}
}

func TestAuthProfiles_LocalProfileBeatsGlobalActiveProfile(t *testing.T) {
	server := newAuthServer(t, []fakeHoneycombAccount{
		{key: "work-key", id: "work-id", keyType: "configuration", team: "work-team", env: "prod", teamName: "Work", envName: "Production"},
		{key: "personal-key", id: "personal-id", keyType: "configuration", team: "personal-team", env: "playground", teamName: "Personal", envName: "Playground"},
	})
	defer server.Close()

	env := isolatedConfigEnv(t)
	projectDir := t.TempDir()

	_, stderr, code := runCLIWithOptions(t, cliRunOptions{Env: env, Stdin: "work-key\n"},
		"--api-url", server.URL,
		"auth", "login", "--profile", "work", "--api-key-stdin", "--no-switch",
	)
	if code != 0 {
		t.Fatalf("login work failed with code %d\nstderr: %s", code, stderr)
	}
	_, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env, Stdin: "personal-key\n"},
		"--api-url", server.URL,
		"auth", "login", "--profile", "personal", "--api-key-stdin",
	)
	if code != 0 {
		t.Fatalf("login personal failed with code %d\nstderr: %s", code, stderr)
	}

	stdout, stderr, code := runCLIWithOptions(t, cliRunOptions{Env: env, Dir: projectDir}, "auth", "switch", "work", "--local")
	if code != 0 {
		t.Fatalf("local switch failed with code %d\nstderr: %s", code, stderr)
	}
	switchOut := parseJSON(t, stdout)
	if switchOut["local_active"] != true {
		t.Fatalf("expected local switch output, got: %s", stdout)
	}

	stdout, stderr, code = runCLIWithOptions(t, cliRunOptions{Env: env, Dir: projectDir}, "auth", "whoami")
	if code != 0 {
		t.Fatalf("whoami with local profile failed with code %d\nstderr: %s", code, stderr)
	}
	assertAuthTeam(t, stdout, "work-team")
}

func newAuthServer(t *testing.T, accounts []fakeHoneycombAccount) *httptest.Server {
	t.Helper()
	byKey := make(map[string]fakeHoneycombAccount, len(accounts))
	for _, account := range accounts {
		byKey[account.key] = account
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/auth" {
			http.NotFound(w, r)
			return
		}
		account, ok := byKey[r.Header.Get("X-Honeycomb-Team")]
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":   account.id,
			"type": account.keyType,
			"api_key_access": map[string]bool{
				"boards": true,
			},
			"team": map[string]string{
				"name": account.teamName,
				"slug": account.team,
			},
			"environment": map[string]string{
				"name": account.envName,
				"slug": account.env,
			},
		})
	}))
}

func assertAuthTeam(t *testing.T, stdout string, wantTeam string) {
	t.Helper()
	m := parseJSON(t, stdout)
	team, ok := m["team"].(map[string]any)
	if !ok {
		t.Fatalf("expected team object, got: %s", stdout)
	}
	if team["slug"] != wantTeam {
		t.Fatalf("expected team %q, got: %s", wantTeam, stdout)
	}
}

func assertProfileInfo(t *testing.T, value any, wantTeam string, wantEnv string) {
	t.Helper()
	info, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected profile_info object, got: %#v", value)
	}
	if info["team_slug"] != wantTeam || info["environment_slug"] != wantEnv {
		t.Fatalf("expected team/env %s/%s, got: %#v", wantTeam, wantEnv, value)
	}
}
