package command

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/crypto"
	"github.com/jarmo/secrets/user"
	"github.com/jarmo/secrets/secret"
	"github.com/jarmo/secrets/storage"
)

func TestAddUser(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "users")
	setupCode, _, err := user.CreatePending(directory, "bob")
	if err != nil {
		t.Fatal(err)
	}
	vaultPath := filepath.Join(t.TempDir(), "production.json")
	if err := storage.Initialize(vaultPath, []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	if err := addUser(vaultPath, directory, "bob", []byte("vault-password")); err == nil {
		t.Fatal("Pending user was granted access")
	}
	if err := user.Activate(directory, "bob", setupCode, []byte("bob-password")); err != nil {
		t.Fatal(err)
	}
	if err := addUser(vaultPath, directory, "bob", []byte("incorrect")); err == nil {
		t.Fatal("Incorrect vault password granted access")
	}
	if err := addUser(vaultPath, directory, "bob", []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	_, privateKey, err := user.Authenticate(directory, "bob", []byte("bob-password"))
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.WriteForUser(vaultPath, privateKey, []secret.Secret{secret.New("api", "token")}); err != nil {
		t.Fatal(err)
	}
	if _, err := storage.Read(vaultPath, []byte("vault-password")); err == nil {
		t.Fatal("Original vault password still unlocks the vault")
	}
	entries, err := storage.ReadForUser(vaultPath, privateKey)
	if err != nil || len(entries) != 1 || entries[0].Value != "token" {
		t.Fatalf("User access broken: %v, %v", entries, err)
	}
	before, err := os.ReadFile(vaultPath)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(before, []byte(`"bob"`)) {
		t.Fatal("Vault file exposes user username")
	}
	if err := addUserForUser(vaultPath, directory, "bob", privateKey); err == nil {
		t.Fatal("Duplicate user association succeeded")
	}
	after, err := os.ReadFile(vaultPath)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("Duplicate association modified vault")
	}
}

func TestAddUserMigratesLegacyVault(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "users")
	if err := user.Register(directory, "alice", []byte("alice-password")); err != nil {
		t.Fatal(err)
	}
	vaultPath := filepath.Join(t.TempDir(), "legacy.json")
	contents := []secret.Secret{secret.New("entry", "value")}
	encrypted := crypto.Encrypt([]byte("legacy-password"), contents)
	data, err := json.Marshal(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(vaultPath, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := addUser(vaultPath, directory, "alice", []byte("bad-password")); err == nil {
		t.Fatal("Incorrect password migrated legacy vault")
	}
	original, _ := os.ReadFile(vaultPath)
	if !bytes.Equal(original, data) {
		t.Fatal("Failed migration modified legacy vault")
	}
	if err := addUser(vaultPath, directory, "alice", []byte("legacy-password")); err != nil {
		t.Fatal(err)
	}
	_, privateKey, err := user.Authenticate(directory, "alice", []byte("alice-password"))
	if err != nil {
		t.Fatal(err)
	}
	entries, err := storage.ReadForUser(vaultPath, privateKey)
	if err != nil || len(entries) != 1 || entries[0] != contents[0] {
		t.Fatalf("Legacy vault contents lost: %v, %v", entries, err)
	}
	if _, err := storage.Read(vaultPath, []byte("legacy-password")); err == nil {
		t.Fatal("Legacy password still opens converted vault")
	}
}
