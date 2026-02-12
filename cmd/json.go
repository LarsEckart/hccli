package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// outputSizeThreshold is the size in bytes above which output is considered large.
// When exceeded, the full output is also written to a temp file and a hint is
// printed to stderr so that callers that truncate stdout can still find it.
//
// Known truncation limits of AI coding agents:
//   - Codex CLI: ~40KB default (10k output tokens × ~4 bytes/token)
//   - Amp:       ~30-50KB (persisted output limit)
//   - pi:        50KB / 2000 lines
//
// 30KB fires before any of these truncate.
const outputSizeThreshold = 30 * 1024 // 30KB

func printJSON(v any) error {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	buf = append(buf, '\n')

	if len(buf) > outputSizeThreshold {
		tmpPath, writeErr := writeTempFile(buf)
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Output is large (%s) and could not be saved to a temp file: %v\n", formatSize(len(buf)), writeErr)
		} else {
			fmt.Fprintf(os.Stderr, "⚠️  Output is large (%s). Full output written to: %s\n", formatSize(len(buf)), tmpPath)
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "💡 To reduce output size:")
			fmt.Fprintln(os.Stderr, "  • Use fewer --breakdown flags")
			fmt.Fprintln(os.Stderr, "  • Use a shorter --time-range")
			fmt.Fprintln(os.Stderr, "  • Add filters to narrow results")
		}
	}

	_, err = os.Stdout.Write(buf)
	return err
}

func writeTempFile(data []byte) (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	name := fmt.Sprintf("hccli-%s.json", hex.EncodeToString(b))
	path := filepath.Join(os.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func formatSize(bytes int) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	}
}
