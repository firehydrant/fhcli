package cli

import (
	"os"
	"testing"
)

// TestMain points XDG_CONFIG_HOME at an empty temp dir so a developer's real
// ~/.config/fhcli/spec.json can't change which spec the tests see.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "fhcli-test-config")
	if err != nil {
		panic(err)
	}
	os.Setenv("XDG_CONFIG_HOME", dir)

	code := m.Run()

	os.RemoveAll(dir)
	os.Exit(code)
}
