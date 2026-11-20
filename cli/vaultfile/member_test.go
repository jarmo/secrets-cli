package vaultfile

import (
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/member"
	"github.com/jarmo/secrets/secret"
	"github.com/jarmo/secrets/storage"
)

func TestOpenWithMemberCredentialsReadsAndWrites(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "members")
	if err := member.Register(directory, "alice", []byte("alice-password")); err != nil {
		t.Fatal(err)
	}
	vaultPath := filepath.Join(t.TempDir(), "vault.json")
	if err := storage.Initialize(vaultPath, []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	alice, err := member.Load(directory, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.AssociateMember(vaultPath, []byte("vault-password"), alice.Identity); err != nil {
		t.Fatal(err)
	}
	openedVault, err := OpenWithPassword(vaultPath, "alice", directory, []byte("alice-password"))
	if err != nil || openedVault.PrivateKey() == nil {
		t.Fatalf("Member cannot open vault: %v", err)
	}
	if err := openedVault.Write([]secret.Secret{secret.New("service", "token")}); err != nil {
		t.Fatal(err)
	}
	openedAgain, err := OpenWithPassword(vaultPath, "alice", directory, []byte("alice-password"))
	if err != nil || len(openedAgain.Secrets) != 1 || openedAgain.Secrets[0].Value != "token" {
		t.Fatalf("Member cannot read updated vault: %v", err)
	}
	legacy, err := OpenWithPassword(vaultPath, "", directory, []byte("vault-password"))
	if err != nil || len(legacy.Secrets) != 1 {
		t.Fatalf("Legacy password access broke: %v", err)
	}
	if _, err := OpenWithPassword(vaultPath, "alice", directory, []byte("wrong-password")); err == nil {
		t.Fatal("Incorrect member password was accepted")
	}
	if err := member.Register(directory, "bob", []byte("bob-password")); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenWithPassword(vaultPath, "bob", directory, []byte("bob-password")); err == nil {
		t.Fatal("Member without vault access was accepted")
	}
}
