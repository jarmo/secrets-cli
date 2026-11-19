package command

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "path/filepath"
  "testing"

  "github.com/jarmo/secrets/storage"
)

func TestInitializeVault(t *testing.T) {
  path := filepath.Join(t.TempDir(), "vault.json")
  if err := initializeVault(path, []byte("password"), []byte("different")); err == nil {
    t.Fatal("Mismatched passwords were accepted")
  }
  if _, err := os.Stat(path); !os.IsNotExist(err) {
    t.Fatalf("Vault was created despite mismatched passwords: %v", err)
  }
  if err := initializeVault(path, []byte("password"), []byte("password")); err != nil {
    t.Fatal(err)
  }
  data, err := ioutil.ReadFile(path)
  if err != nil {
    t.Fatal(err)
  }
  var header struct{ Version int }
  if err := json.Unmarshal(data, &header); err != nil || header.Version != 2 {
    t.Fatalf("Initialized vault is not v2: %v", err)
  }
  entries, err := storage.Read(path, []byte("password"))
  if err != nil || len(entries) != 0 {
    t.Fatalf("Newly initialized vault cannot be opened: %v", err)
  }
  if err := initializeVault(path, []byte("password"), []byte("password")); !os.IsExist(err) {
    t.Fatalf("Existing vault was overwritten: %v", err)
  }
}

