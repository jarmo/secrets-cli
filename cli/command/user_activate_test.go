package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jarmo/secrets/user"
)

func TestActivateUser(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "users")
	code, _, err := user.CreatePending(directory, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if err := activateUser(directory, "bob", code, []byte("first"), []byte("different")); err == nil {
		t.Fatal("Accepted mismatched passwords")
	}
	if _, err := user.Load(directory, "bob"); !os.IsNotExist(err) {
		t.Fatalf("Mismatched passwords created an account: %v", err)
	}
	if err := activateUser(directory, "bob", "incorrect", []byte("correct"), []byte("correct")); err == nil {
		t.Fatal("Accepted incorrect setup code")
	}
	if err := activateUser(directory, "bob", code, []byte("correct"), []byte("correct")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := user.Authenticate(directory, "bob", []byte("correct")); err != nil {
		t.Fatal(err)
	}
	if err := activateUser(directory, "bob", code, []byte("second"), []byte("second")); err == nil {
		t.Fatal("Reused setup code")
	}
}
