package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/member"
)

func TestActivateMember(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "members")
	code, _, err := member.CreatePending(directory, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if err := activateMember(directory, "bob", code, []byte("first"), []byte("different")); err == nil {
		t.Fatal("Accepted mismatched passwords")
	}
	if _, err := member.Load(directory, "bob"); !os.IsNotExist(err) {
		t.Fatalf("Mismatched passwords created an account: %v", err)
	}
	if err := activateMember(directory, "bob", "incorrect", []byte("correct"), []byte("correct")); err == nil {
		t.Fatal("Accepted incorrect setup code")
	}
	if err := activateMember(directory, "bob", code, []byte("correct"), []byte("correct")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := member.Authenticate(directory, "bob", []byte("correct")); err != nil {
		t.Fatal(err)
	}
	if err := activateMember(directory, "bob", code, []byte("second"), []byte("second")); err == nil {
		t.Fatal("Reused setup code")
	}
}
