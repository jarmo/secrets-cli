package command

import (
	"os"
	"path/filepath"
	"strings"
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
	if err := activateUser(directory, "bob", "incorrect", []byte("correct-password"), []byte("correct-password")); err == nil {
		t.Fatal("Accepted incorrect setup code")
	}
	for _, candidate := range []string{"short", strings.Repeat("x", 1025)} {
		if err := activateUser(directory, "bob", code, []byte(candidate), []byte(candidate)); err == nil {
			t.Fatal("Accepted invalid password length")
		}
		if _, err := user.Load(directory, "bob"); !os.IsNotExist(err) {
			t.Fatalf("Invalid password created an account: %v", err)
		}
	}
	if err := activateUser(directory, "bob", code, []byte("correct-password"), []byte("correct-password")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := user.Authenticate(directory, "bob", []byte("correct-password")); err != nil {
		t.Fatal(err)
	}
	if err := activateUser(directory, "bob", code, []byte("second"), []byte("second")); err == nil {
		t.Fatal("Reused setup code")
	}
}
