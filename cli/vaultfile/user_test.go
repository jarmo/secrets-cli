package vaultfile

import (
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/user"
	"github.com/jarmo/secrets/secret"
	"github.com/jarmo/secrets/storage"
)

func TestOpenWithUserCredentialsReadsAndWrites(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "users")
	if err := user.Register(directory, "alice", []byte("alice-password")); err != nil {
		t.Fatal(err)
	}
	vaultPath := filepath.Join(t.TempDir(), "vault.json")
	if err := storage.Initialize(vaultPath, []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	alice, err := user.Load(directory, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.AssociateUser(vaultPath, []byte("vault-password"), alice.Identity); err != nil {
		t.Fatal(err)
	}
	openedVault, err := OpenWithPassword(vaultPath, "alice", directory, []byte("alice-password"))
	if err != nil || openedVault.PrivateKey() == nil {
		t.Fatalf("User cannot open vault: %v", err)
	}
	if err := openedVault.Write([]secret.Secret{secret.New("service", "token")}); err != nil {
		t.Fatal(err)
	}
	openedAgain, err := OpenWithPassword(vaultPath, "alice", directory, []byte("alice-password"))
	if err != nil || len(openedAgain.Secrets) != 1 || openedAgain.Secrets[0].Value != "token" {
		t.Fatalf("User cannot read updated vault: %v", err)
	}
	legacy, err := OpenWithPassword(vaultPath, "", directory, []byte("vault-password"))
	if err != nil || len(legacy.Secrets) != 1 {
		t.Fatalf("Legacy password access broke: %v", err)
	}
	if _, err := OpenWithPassword(vaultPath, "alice", directory, []byte("wrong-password")); err == nil {
		t.Fatal("Incorrect user password was accepted")
	}
	if err := user.Register(directory, "bob", []byte("bob-password")); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenWithPassword(vaultPath, "bob", directory, []byte("bob-password")); err == nil {
		t.Fatal("User without vault access was accepted")
	}
}
