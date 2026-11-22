package path

import (
  "os"
  "path/filepath"
  "testing"
)

func TestUsersDirectoryPrefersConfiguration(t *testing.T) {
  configHome := t.TempDir()
  t.Setenv("XDG_CONFIG_HOME", configHome)
  directory := filepath.Join(configHome, "secrets")
  if err := os.MkdirAll(directory, 0700); err != nil {
    t.Fatal(err)
  }
  if expected := filepath.Join(directory, "users"); UsersDirectory() != expected {
    t.Fatalf("Unexpected default directory: %q", UsersDirectory())
  }
  configured := filepath.Join(t.TempDir(), "Dropbox", "users")
  contents := `{"UsersDir":"` + filepath.ToSlash(configured) + `","Vaults":[]}`
  if err := os.WriteFile(filepath.Join(directory, "config.json"), []byte(contents), 0600); err != nil {
    t.Fatal(err)
  }
  if actual := UsersDirectory(); actual != configured {
    t.Fatalf("Unexpected configured directory: %q", actual)
  }
}
