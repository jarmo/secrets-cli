package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets-cli/v7/storage/path"
	"github.com/jarmo/secrets/user"
	"github.com/jarmo/secrets/storage"
)

func TestUserDestroyRevokesConfiguredVaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	directory := path.UsersDirectory()
	for _, username := range []string{"alice", "bob"} {
		if err := user.Register(directory, username, []byte(username+"-password")); err != nil {
			t.Fatal(err)
		}
	}
	alice, aliceKey, _ := user.Authenticate(directory, "alice", []byte("alice-password"))
	bob, bobKey, _ := user.Authenticate(directory, "bob", []byte("bob-password"))
	vaultPath := filepath.Join(t.TempDir(), "vault.json")
	if err := storage.InitializeForUser(vaultPath, alice.Identity); err != nil {
		t.Fatal(err)
	}
	if err := storage.AssociateUserForUser(vaultPath, aliceKey, bob.Identity); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path.UsersDirectory()), 0700); err != nil {
		t.Fatal(err)
	}
	path.Store(vaultPath, "production")
	command := UserDestroy{Username: "bob", User: "alice", UsersDirectory: directory}
	if err := command.destroy([]byte("wrong-password")); err == nil {
		t.Fatal("Incorrect user password destroyed the account")
	}
	if err := command.destroy([]byte("alice-password")); err != nil {
		t.Fatal(err)
	}
	if _, err := user.Load(directory, "bob"); !os.IsNotExist(err) {
		t.Fatalf("User account still exists: %v", err)
	}
	if _, err := storage.ReadForUser(vaultPath, bobKey); err == nil {
		t.Fatal("Destroyed user can still decrypt the vault")
	}
	if _, err := storage.ReadForUser(vaultPath, aliceKey); err != nil {
		t.Fatal(err)
	}
}
