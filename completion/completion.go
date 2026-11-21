package completion

import (
	"embed"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/jarmo/secrets-cli/v7/storage/path"
	"github.com/jarmo/secrets/user"
)

//go:embed bash zsh
var scripts embed.FS

func Script(output io.Writer, shell string) error {
	if shell != "bash" && shell != "zsh" {
		return fmt.Errorf("Unsupported shell: %s", shell)
	}
	contents, err := scripts.ReadFile(shell)
	if err != nil {
		return err
	}
	_, err = output.Write(contents)
	return err
}

func Candidates(output io.Writer, kind, directory string) {
	var candidates []string
	if kind == "aliases" {
		candidates, _ = path.Aliases()
	} else if kind == "users" || kind == "pending" {
		if directory == "" {
			directory = path.UsersDirectory()
		}
		entries, err := os.ReadDir(directory)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if !entry.Type().IsRegular() {
				continue
			}
			if kind == "pending" && strings.HasSuffix(entry.Name(), ".pending.json") {
				username := strings.TrimSuffix(entry.Name(), ".pending.json")
				if validUsername(username) {
					candidates = append(candidates, username)
				}
			} else if kind == "users" && strings.HasSuffix(entry.Name(), ".json") && !strings.HasSuffix(entry.Name(), ".pending.json") {
				username := strings.TrimSuffix(entry.Name(), ".json")
				if _, err := user.Load(directory, username); err == nil {
					candidates = append(candidates, username)
				}
			}
		}
	}
	sort.Strings(candidates)
	for _, candidate := range candidates {
		if candidate != "" && !strings.ContainsAny(candidate, "\r\n") {
			fmt.Fprintln(output, candidate)
		}
	}
}

func validUsername(username string) bool {
	if len(username) == 0 || len(username) > 64 {
		return false
	}
	for index, char := range username {
		if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' {
			continue
		}
		if index > 0 && (char == '.' || char == '_' || char == '-') {
			continue
		}
		return false
	}
	return true
}
