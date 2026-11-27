package command

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "path/filepath"
  "testing"

  "github.com/jarmo/secrets/user"
  "github.com/jarmo/secrets/storage"
  "github.com/jarmo/secrets-cli/v7/storage/path"
)

func TestInitializeVault(t *testing.T) {
  directory := filepath.Join(t.TempDir(), "users")
  vaultPath := filepath.Join(t.TempDir(), "vault.json")
  if err := initializeVault(vaultPath, directory, "alice"); err == nil {
    t.Fatal("Missing user was accepted")
  }
  setupCode, _, err := user.CreatePending(directory, "alice")
  if err != nil {
    t.Fatal(err)
  }
  if err := initializeVault(vaultPath, directory, "alice"); err == nil {
    t.Fatal("Pending user was accepted")
  }
  if _, err := os.Stat(vaultPath); !os.IsNotExist(err) {
    t.Fatalf("Vault exists for missing or pending user: %v", err)
  }
  if err := user.Activate(directory, "alice", setupCode, []byte("alice-password")); err != nil {
    t.Fatal(err)
  }
  if err := initializeVault(vaultPath, directory, "alice"); err != nil {
    t.Fatal(err)
  }
  data, err := ioutil.ReadFile(vaultPath)
  if err != nil {
    t.Fatal(err)
  }
  var header struct{ Version int }
  if err := json.Unmarshal(data, &header); err != nil || header.Version != 2 {
    t.Fatalf("Initialized vault is not v2: %v", err)
  }
  _, privateKey, err := user.Authenticate(directory, "alice", []byte("alice-password"))
  if err != nil {
    t.Fatal(err)
  }
  entries, err := storage.ReadForUser(vaultPath, privateKey)
  if err != nil || len(entries) != 0 {
    t.Fatalf("New user cannot open initialized vault: %v", err)
  }
  if _, err := storage.Read(vaultPath, []byte("alice-password")); err == nil {
    t.Fatal("User-only vault unexpectedly allows password-based access")
  }
  if err := user.Register(directory, "bob", []byte("bob-password")); err != nil {
    t.Fatal(err)
  }
  _, bobKey, err := user.Authenticate(directory, "bob", []byte("bob-password"))
  if err != nil {
    t.Fatal(err)
  }
  if _, err := storage.ReadForUser(vaultPath, bobKey); err == nil {
    t.Fatal("Another user can open the vault")
  }
  if err := initializeVault(vaultPath, directory, "alice"); !os.IsExist(err) {
    t.Fatalf("Existing vault was overwritten: %v", err)
  }
}

func TestInitializeExecuteUsesConfiguredUsersDirectory(t *testing.T) {
  configHome := t.TempDir()
  t.Setenv("XDG_CONFIG_HOME", configHome)
  configDirectory := filepath.Join(configHome, "secrets")
  if err := os.MkdirAll(configDirectory, 0700); err != nil {
    t.Fatal(err)
  }
  directory := filepath.Join(t.TempDir(), "users")
  if err := user.Register(directory, "alice", []byte("alice-password")); err != nil {
    t.Fatal(err)
  }
  config, err := json.Marshal(map[string]interface{}{"UsersDir": directory, "Vaults": []interface{}{}})
  if err != nil {
    t.Fatal(err)
  }
  if err := os.WriteFile(filepath.Join(configDirectory, "config.json"), config, 0600); err != nil {
    t.Fatal(err)
  }
  vaultPath := filepath.Join(t.TempDir(), "vault.json")
  Initialize{VaultPath: vaultPath, VaultAlias: "personal", User: "alice"}.Execute()
  if configuredPath, err := path.Get("personal"); err != nil || configuredPath != vaultPath {
    t.Fatalf("Vault alias was not configured: %q, %v", configuredPath, err)
  }
  if path.UsersDirectory() != directory {
    t.Fatal("Configured users directory changed")
  }
  _, privateKey, err := user.Authenticate(directory, "alice", []byte("alice-password"))
  if err != nil {
    t.Fatal(err)
  }
  if _, err := storage.ReadForUser(vaultPath, privateKey); err != nil {
    t.Fatal(err)
  }
}
