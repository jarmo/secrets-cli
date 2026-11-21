package completion

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestShellCompletions(t *testing.T) {
	binDirectory := t.TempDir()
	binary := filepath.Join(binDirectory, "secrets")
	stub := `#!/bin/sh
if [ "$1" = __complete ]; then
  case "$2" in
    aliases) printf 'development\nproduction\n' ;;
    users) printf 'alice\nbob\n' ;;
    pending) printf 'charlie\n' ;;
  esac
fi
`
	if err := os.WriteFile(binary, []byte(stub), 0700); err != nil {
		t.Fatal(err)
	}
	zshScript, err := os.ReadFile("zsh")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(binDirectory, "_secrets"), zshScript, 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		words  string
		want   string
		reject string
	}{
		{"command names", "secrets ''", "user-add", "alice"},
		{"global version", "secrets ''", "--version", "alice"},
		{"list options", "secrets list ''", "--alias", "alice"},
		{"list options after filter", "secrets list search ''", "--alias", "alice"},
		{"add options after name", "secrets add secret ''", "--alias", "alice"},
		{"edit options after ID", "secrets edit id ''", "--alias", "alice"},
		{"delete options after ID", "secrets delete id ''", "--alias", "alice"},
		{"initialize options", "secrets initialize ''", "--user", "alice"},
		{"user-create options after username", "secrets user-create new ''", "--users-dir", "alice"},
		{"user-activate pending username", "secrets user-activate ''", "charlie", "--users-dir"},
		{"user-activate options after username", "secrets user-activate charlie ''", "--users-dir", "alice"},
		{"user-add username", "secrets user-add ''", "alice", "--alias"},
		{"user-add options after username", "secrets user-add bob ''", "--alias", "alice"},
		{"user-list options", "secrets user-list ''", "--alias", "alice"},
		{"user-remove username", "secrets user-remove ''", "alice", "--alias"},
		{"user-remove options after username", "secrets user-remove bob ''", "--alias", "alice"},
		{"user-destroy username", "secrets user-destroy ''", "alice", "--user"},
		{"user-destroy options after username", "secrets user-destroy bob ''", "--user", "alice"},
		{"user-password-change username", "secrets user-password-change ''", "alice", "--users-dir"},
		{"user-password-change options after username", "secrets user-password-change bob ''", "--users-dir", "alice"},
		{"completion shells", "secrets completion ''", "zsh", "--alias"},
		{"user-add after alias option", "secrets user-add --alias production ''", "alice", "--alias"},
		{"user-add after inline alias", "secrets user-add --alias=production ''", "alice", "--alias"},
		{"user-remove after path option", "secrets user-remove --path vault.json ''", "alice", "--alias"},
		{"user-destroy after acting user", "secrets user-destroy --user alice ''", "bob", "--user"},
		{"user-password-change after users directory", "secrets user-password-change --users-dir directory ''", "alice", "--users-dir"},
		{"user-activate after users directory", "secrets user-activate --users-dir directory ''", "charlie", "--users-dir"},
		{"completion after username and alias", "secrets user-add --alias production bob ''", "--alias", "alice"},
		{"alias option value", "secrets list --alias ''", "production", "alice"},
		{"user option value", "secrets list --user ''", "alice", "--alias"},
		{"list option after inline path", "secrets list --path=vault.json ''", "--alias", "alice"},
		{"user-add option after user", "secrets user-add bob --user alice ''", "--alias", "bob"},
	}

	for _, shell := range []string{"bash", "zsh"} {
		t.Run(shell, func(t *testing.T) {
			if _, err := exec.LookPath(shell); err != nil {
				t.Skipf("%s is not installed", shell)
			}
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					position := len(strings.Fields(test.words))
					var program string
					if shell == "bash" {
						program = "source ./bash\nCOMP_WORDS=(" + test.words + ")\nCOMP_CWORD=" + strconv.Itoa(position-1) + "\n_secrets_completion\nprintf '%s\\n' \"${COMPREPLY[@]}\"\n"
					} else {
						program = "fpath=(" + strconv.Quote(binDirectory) + " $fpath)\nautoload -Uz _secrets\ncompadd() { while (( $# )); do if [[ \"$1\" == -- ]]; then shift; break; fi; shift; done; print -rl -- \"$@\"; }\nwords=(" + test.words + ")\nCURRENT=" + strconv.Itoa(position) + "\n_secrets\n"
					}
					command := exec.Command(shell, "-c", program)
					command.Env = append(os.Environ(), "PATH="+binDirectory+string(os.PathListSeparator)+os.Getenv("PATH"))
					output, err := command.CombinedOutput()
					if err != nil {
						t.Fatalf("Completion failed: %s: %v", output, err)
					}
					candidates := strings.Fields(string(output))
					if !contains(candidates, test.want) || contains(candidates, test.reject) {
						t.Fatalf("Unexpected completion candidates: %v", candidates)
					}
				})
			}
		})
	}
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func TestZshCompletionAutoloadFormat(t *testing.T) {
	script, err := os.ReadFile("zsh")
	if err != nil {
		t.Fatal(err)
	}
	contents := string(script)
	if !strings.HasPrefix(contents, "#compdef secrets\n") {
		t.Fatal("Zsh completion must declare its command for compinit")
	}
	if strings.Contains(contents, "_secrets() {") || strings.Contains(contents, "compdef _secrets secrets") {
		t.Fatal("Zsh completion must contain an autoload function body, not a nested function definition")
	}
}

func TestZshCompinitRegistersCommand(t *testing.T) {
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("zsh is not installed")
	}
	completionDirectory := t.TempDir()
	script, err := os.ReadFile("zsh")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(completionDirectory, "_secrets"), script, 0600); err != nil {
		t.Fatal(err)
	}
	program := "fpath=(" + strconv.Quote(completionDirectory) + " $fpath)\nautoload -Uz compinit\ncompinit -i -D\n[[ \"${_comps[secrets]}\" == _secrets ]]\n"
	output, err := exec.Command("zsh", "-fc", program).CombinedOutput()
	if err != nil {
		t.Fatalf("Zsh compinit did not register secrets: %v\n%s", err, output)
	}
}
