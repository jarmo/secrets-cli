package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v7/cli/vaultfile"
	"github.com/jarmo/secrets-cli/v7/input"
	"github.com/jarmo/secrets-cli/v7/storage/path"
	"github.com/jarmo/secrets/storage"
)

type UserRemove struct {
	Username         string
	VaultAlias       string
	VaultPath        string
	UsersDirectory string
	User           string
}

func (command UserRemove) Execute() {
	directory := command.UsersDirectory
	if directory == "" {
		directory = path.UsersDirectory()
	}
	vaultPath := vaultfile.Path(command.VaultAlias, command.VaultPath)
	var err error
	if command.User != "" {
		openedVault := vaultfile.Open(vaultPath, command.User, directory)
		err = storage.RemoveUserForUser(vaultPath, directory, command.Username, openedVault.PrivateKey())
	} else {
		err = storage.RemoveUser(vaultPath, directory, command.Username, input.AskVaultPassword())
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("User %s removed from vault!\n", command.Username)
}
