package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/secret"
	"github.com/jarmo/secrets/storage"
	"github.com/jarmo/secrets/user"
)

func TestUserPasswordChangePreservesVaultAccess(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "users")
	oldPassword := []byte("alice-old-password")
	newPassword := []byte("alice-new-password")
	if err := user.Register(directory, "alice", oldPassword); err != nil {
		t.Fatal(err)
	}
	account, privateKey, err := user.Authenticate(directory, "alice", oldPassword)
	if err != nil {
		t.Fatal(err)
	}
	vaultPaths := []string{
		filepath.Join(t.TempDir(), "production.json"),
		filepath.Join(t.TempDir(), "development.json"),
	}
	for _, vaultPath := range vaultPaths {
		if err := storage.InitializeForUser(vaultPath, account.Identity); err != nil {
			t.Fatal(err)
		}
		if err := storage.WriteForUser(vaultPath, privateKey, []secret.Secret{secret.New("name", "value")}); err != nil {
			t.Fatal(err)
		}
	}
	before := make([][]byte, len(vaultPaths))
	for i, vaultPath := range vaultPaths {
		before[i], err = os.ReadFile(vaultPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	accountFile := filepath.Join(directory, "alice.json")
	originalAccount, err := os.ReadFile(accountFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, candidate := range []struct {
		name         string
		oldPassword  []byte
		newPassword  []byte
		confirmation []byte
	}{
		{"wrong current password", []byte("wrong-password"), newPassword, newPassword},
		{"confirmation mismatch", oldPassword, newPassword, []byte("does-not-match")},
		{"short password", oldPassword, []byte("short"), []byte("short")},
		{"unchanged password", oldPassword, oldPassword, oldPassword},
	} {
		t.Run(candidate.name, func(t *testing.T) {
			if err := changeUserPassword(directory, "alice", candidate.oldPassword, candidate.newPassword, candidate.confirmation); err == nil {
				t.Fatal("Invalid password change succeeded")
			}
			contents, err := os.ReadFile(accountFile)
			if err != nil || !bytes.Equal(contents, originalAccount) {
				t.Fatalf("Failed change modified account: %v", err)
			}
		})
	}
	if err := changeUserPassword(directory, "alice", oldPassword, newPassword, newPassword); err != nil {
		t.Fatal(err)
	}
	if _, _, err := user.Authenticate(directory, "alice", oldPassword); err == nil {
		t.Fatal("Old password still authenticates")
	}
	_, updatedKey, err := user.Authenticate(directory, "alice", newPassword)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(privateKey.Bytes(), updatedKey.Bytes()) {
		t.Fatal("Password change replaced user identity")
	}
	for i, vaultPath := range vaultPaths {
		contents, err := os.ReadFile(vaultPath)
		if err != nil || !bytes.Equal(contents, before[i]) {
			t.Fatalf("Password change modified vault %s: %v", vaultPath, err)
		}
		entries, err := storage.ReadForUser(vaultPath, updatedKey)
		if err != nil || len(entries) != 1 || entries[0].Value != "value" {
			t.Fatalf("Cannot open vault %s: %v", vaultPath, err)
		}
	}
}
