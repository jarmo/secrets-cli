package command

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "path/filepath"
  "testing"

  "github.com/jarmo/secrets/member"
  "github.com/jarmo/secrets/storage"
)

func TestInitializeVault(t *testing.T) {
  directory := filepath.Join(t.TempDir(), "members")
  vaultPath := filepath.Join(t.TempDir(), "vault.json")
  if err := initializeVault(vaultPath, directory, "alice", []byte("password")); err == nil {
    t.Fatal("Missing member was accepted")
  }
  setupCode, _, err := member.CreatePending(directory, "alice")
  if err != nil {
    t.Fatal(err)
  }
  if err := initializeVault(vaultPath, directory, "alice", []byte("password")); err == nil {
    t.Fatal("Pending member was accepted")
  }
  if err := member.Activate(directory, "alice", setupCode, []byte("alice-password")); err != nil {
    t.Fatal(err)
  }
  if err := initializeVault(vaultPath, directory, "alice", []byte("incorrect")); err == nil {
    t.Fatal("Incorrect member password was accepted")
  }
  if _, err := os.Stat(vaultPath); !os.IsNotExist(err) {
    t.Fatalf("Vault exists despite failed authentication: %v", err)
  }
  if err := initializeVault(vaultPath, directory, "alice", []byte("alice-password")); err != nil {
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
  _, privateKey, err := member.Authenticate(directory, "alice", []byte("alice-password"))
  if err != nil {
    t.Fatal(err)
  }
  entries, err := storage.ReadForMember(vaultPath, privateKey)
  if err != nil || len(entries) != 0 {
    t.Fatalf("New member cannot open initialized vault: %v", err)
  }
  if _, err := storage.Read(vaultPath, []byte("alice-password")); err == nil {
    t.Fatal("Member-only vault unexpectedly allows password-based access")
  }
  if err := initializeVault(vaultPath, directory, "alice", []byte("alice-password")); !os.IsExist(err) {
    t.Fatalf("Existing vault was overwritten: %v", err)
  }
}
