package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/cli/vaultfile"
	"github.com/jarmo/secrets-cli/v6/input"
	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/crypto"
	"github.com/jarmo/secrets/storage"
)

type UserList struct {
	VaultAlias       string
	VaultPath        string
	UsersDirectory string
	User           string
}

func (command UserList) Execute() {
	directory := command.UsersDirectory
	if directory == "" {
		directory = path.UsersDirectory()
	}
	vaultPath := vaultfile.Path(command.VaultAlias, command.VaultPath)
	var users []crypto.VaultUser
	var err error
	if command.User != "" {
		openedVault := vaultfile.Open(vaultPath, command.User, directory)
		users, err = storage.ListUsersForUser(vaultPath, directory, openedVault.PrivateKey())
	} else {
		users, err = storage.ListUsers(vaultPath, directory, input.AskVaultPassword())
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, user := range users {
		if user.PasswordBased {
			fmt.Printf("%s (vault password)\n", user.Name)
		} else {
			fmt.Println(user.Name)
		}
	}
}
