package completion

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarmo/secrets/user"
)

func TestScripts(t *testing.T) {
	for _, shell := range []string{"bash", "zsh"} {
		var output bytes.Buffer
		if err := Script(&output, shell); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(output.String(), "secrets __complete aliases") || !strings.Contains(output.String(), "user-password-change") {
			t.Fatalf("Incomplete %s script", shell)
		}
	}
	if err := Script(&bytes.Buffer{}, "fish"); err == nil {
		t.Fatal("Expected unsupported shell to fail")
	}
}

func TestAliasCandidates(t *testing.T) {
	configurationRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configurationRoot)
	configurationDirectory := filepath.Join(configurationRoot, "secrets")
	if err := os.MkdirAll(configurationDirectory, 0700); err != nil {
		t.Fatal(err)
	}
	configuration := `[{"Path":"/tmp/z","Alias":"zeta"},{"Path":"/tmp/a","Alias":"alpha"}]`
	if err := os.WriteFile(filepath.Join(configurationDirectory, "config.json"), []byte(configuration), 0600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	Candidates(&output, "aliases", "")
	if output.String() != "alpha\nzeta\n" {
		t.Fatalf("Unexpected aliases: %q", output.String())
	}
}

func TestUserCandidates(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "users")
	if err := user.Register(directory, "zoe", []byte("a-long-unique-password")); err != nil {
		t.Fatal(err)
	}
	if err := user.Register(directory, "alice", []byte("another-unique-password")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := user.CreatePending(directory, "pending"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "corrupt.json"), []byte("bad"), 0600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	Candidates(&output, "users", directory)
	if output.String() != "alice\nzoe\n" {
		t.Fatalf("Unexpected active users: %q", output.String())
	}
	output.Reset()
	Candidates(&output, "pending", directory)
	if output.String() != "pending\n" {
		t.Fatalf("Unexpected pending users: %q", output.String())
	}
	output.Reset()
	Candidates(&output, "users", filepath.Join(directory, "missing"))
	if output.Len() != 0 {
		t.Fatalf("Expected empty completion for a missing directory: %q", output.String())
	}
}
