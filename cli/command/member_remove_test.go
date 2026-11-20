package command

import (
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/member"
	"github.com/jarmo/secrets/storage"
)

func TestMemberListAndRemove(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "members")
	if err := member.Register(directory, "bob", []byte("bob-password")); err != nil {
		t.Fatal(err)
	}
	record, privateKey, err := member.Authenticate(directory, "bob", []byte("bob-password"))
	if err != nil {
		t.Fatal(err)
	}
	vaultPath := filepath.Join(t.TempDir(), "vault.json")
	if err := storage.Initialize(vaultPath, []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	if err := storage.AssociateMember(vaultPath, []byte("vault-password"), record.Identity); err != nil {
		t.Fatal(err)
	}
	members, err := storage.ListMembers(vaultPath, directory, []byte("vault-password"))
	if err != nil || len(members) != 2 || members[1].Name != "bob" {
		t.Fatalf("Expected member bob: %+v, %v", members, err)
	}
	if err := storage.RemoveMember(vaultPath, directory, "bob", []byte("vault-password")); err != nil {
		t.Fatal(err)
	}
	if _, err := storage.ReadForMember(vaultPath, privateKey); err == nil {
		t.Fatal("Removed member can still access vault")
	}
}
