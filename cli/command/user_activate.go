package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v7/input"
	"github.com/jarmo/secrets-cli/v7/storage/path"
	"github.com/jarmo/secrets/user"
)

type UserActivate struct {
	Username         string
	UsersDirectory string
}

func (command UserActivate) Execute() {
	directory := command.UsersDirectory
	if directory == "" {
		directory = path.UsersDirectory()
	}
	code := input.AskPassword("Enter setup code: ")
	password := input.AskPassword("Enter new user password: ")
	confirmation := input.AskPassword("Enter new user password again: ")
	if err := activateUser(directory, command.Username, string(code), password, confirmation); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("User %s activated!\n", command.Username)
}

func activateUser(directory, username, code string, password, confirmation []byte) error {
	if !bytes.Equal(password, confirmation) {
		return errors.New("Passwords do not match")
	}
	return user.Activate(directory, username, code, password)
}
