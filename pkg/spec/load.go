package spec

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// embedded is a snapshot of https://api.firehydrant.io/v1/swagger_doc.
// Refresh it with `make spec`.
//
//go:embed spec.json
var embedded []byte

// EmbeddedSource is Spec.Source for the spec compiled into the binary.
const EmbeddedSource = "embedded"

// OverridePath is where `fhcli spec update` writes a newer spec.
func OverridePath() (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "fhcli", "spec.json"), nil
}

// Load returns the override spec if one exists, else the embedded snapshot.
func Load() (*Spec, error) {
	path, err := OverridePath()
	if err != nil {
		return loadEmbedded()
	}

	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return loadEmbedded()
	}
	if err != nil {
		return nil, fmt.Errorf("read spec override %s: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read spec override %s: %w", path, err)
	}
	s, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("spec override %s is invalid (run `fhcli spec update` or delete the file): %w", path, err)
	}
	s.Source = path
	s.ModifiedAt = info.ModTime()
	return s, nil
}

func loadEmbedded() (*Spec, error) {
	s, err := Parse(embedded)
	if err != nil {
		return nil, fmt.Errorf("embedded spec: %w", err)
	}
	s.Source = EmbeddedSource
	return s, nil
}
