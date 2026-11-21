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

type UserPasswordChange struct {
	Username       string
	UsersDirectory string
}

func (command UserPasswordChange) Execute() {
	if err := command.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("User password successfully changed!")
}

func (command UserPasswordChange) run() error {
	currentPassword := input.AskPassword("Enter current user password: ")
	newPassword := input.AskPassword("Enter new user password: ")
	confirmation := input.AskPassword("Enter new user password again: ")
	defer func() {
		for _, password := range [][]byte{currentPassword, newPassword, confirmation} {
			for i := range password {
				password[i] = 0
			}
		}
	}()

	directory := command.UsersDirectory
	if directory == "" {
		directory = path.UsersDirectory()
	}
	return changeUserPassword(directory, command.Username, currentPassword, newPassword, confirmation)
}

func changeUserPassword(directory, username string, currentPassword, newPassword, confirmation []byte) error {
	if err := user.ValidatePassword(newPassword); err != nil {
		return err
	}
	if !bytes.Equal(newPassword, confirmation) {
		return errors.New("Passwords do not match")
	}
	if bytes.Equal(currentPassword, newPassword) {
		return errors.New("Choose a different password")
	}
	return user.ChangePassword(directory, username, currentPassword, newPassword)
}
