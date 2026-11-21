package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v7/input"
	"github.com/jarmo/secrets-cli/v7/storage/path"
	"github.com/jarmo/secrets/user"
	"github.com/jarmo/secrets/storage"
)

type UserDestroy struct {
	Username         string
	User           string
	UsersDirectory string
}

func (command UserDestroy) Execute() {
	if err := command.destroy(input.AskPassword("Enter user password: ")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("User %s destroyed!\n", command.Username)
}

func (command UserDestroy) destroy(password []byte) error {
	defer func() {
		for i := range password {
			password[i] = 0
		}
	}()
	directory := command.UsersDirectory
	if directory == "" {
		directory = path.UsersDirectory()
	}
	_, privateKey, err := user.Authenticate(directory, command.User, password)
	if err != nil {
		return err
	}
	aliases, err := path.Aliases()
	if err != nil {
		return err
	}
	vaultPaths := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		vaultPath, err := path.Get(alias)
		if err != nil {
			return err
		}
		vaultPaths = append(vaultPaths, vaultPath)
	}
	return storage.DestroyUser(vaultPaths, directory, command.Username, privateKey)
}
