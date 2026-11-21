package command

import (
	"bytes"
	"testing"

	"github.com/jarmo/secrets/crypto"
)

func TestPrintUsersSorted(t *testing.T) {
	users := []crypto.VaultUser{
		{Name: "zoe"},
		{Name: "owner", PasswordBased: true},
		{Name: "alice"},
		{Name: "bob"},
	}
	var output bytes.Buffer
	printUsers(&output, users)
	want := "alice\nbob\nowner (vault password)\nzoe\n"
	if output.String() != want {
		t.Fatalf("Unexpected user list:\n%s\nwant:\n%s", output.String(), want)
	}
}
