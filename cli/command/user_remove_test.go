package command

import (
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/user"
	"github.com/jarmo/secrets/storage"
)

func TestUserListAndRemove(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "users")
	if err := user.Register(directory, "bob", []byte("bob-password")); err != nil {
		t.Fatal(err)
	}
	record, privateKey, err := user.Authenticate(directory, "bob", []byte("bob-password"))
	if err != nil {
		t.Fatal(err)
	}
	vaultPath := filepath.Join(t.TempDir(), "vault.json")
	if err := storage.Initialize(vaultPath, []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	if err := storage.AssociateUser(vaultPath, []byte("vault-password"), record.Identity); err != nil {
		t.Fatal(err)
	}
	users, err := storage.ListUsers(vaultPath, directory, []byte("vault-password"))
	if err != nil || len(users) != 2 || users[1].Name != "bob" {
		t.Fatalf("Expected user bob: %+v, %v", users, err)
	}
	if err := storage.RemoveUser(vaultPath, directory, "bob", []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	if _, err := storage.ReadForUser(vaultPath, privateKey); err == nil {
		t.Fatal("Removed user can still access vault")
	}
}
