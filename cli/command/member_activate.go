package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/input"
	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/member"
)

type MemberActivate struct {
	Username         string
	MembersDirectory string
}

func (command MemberActivate) Execute() {
	directory := command.MembersDirectory
	if directory == "" {
		directory = path.MembersDirectory()
	}
	code := input.AskPassword("Enter setup code: ")
	password := input.AskPassword("Enter new member password: ")
	confirmation := input.AskPassword("Enter new member password again: ")
	if err := activateMember(directory, command.Username, string(code), password, confirmation); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Member %s activated!\n", command.Username)
}

func activateMember(directory, username, code string, password, confirmation []byte) error {
	if !bytes.Equal(password, confirmation) {
		return errors.New("Passwords do not match")
	}
	return member.Activate(directory, username, code, password)
}
