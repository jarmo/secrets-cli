package cli

import (
	"testing"

	"github.com/jarmo/secrets-cli/v6/cli/command"
)

func TestCommandCompletion(t *testing.T) {
	for _, shell := range []string{"bash", "zsh"} {
		parsed, ok := Command(version, []string{"completion", shell}).(command.Completion)
		if !ok || parsed.Shell != shell {
			t.Fatalf("Unexpected completion command for %s: %#v", shell, parsed)
		}
	}
	parsed, ok := Command(version, []string{"__complete", "users", "/tmp/user-list"}).(command.Completion)
	if !ok || parsed.CandidateKind != "users" || parsed.UsersDirectory != "/tmp/user-list" {
		t.Fatalf("Unexpected internal completion command: %#v", parsed)
	}
}
