package command

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/jarmo/secrets-cli/v7/cli/vaultfile"
	"github.com/jarmo/secrets-cli/v7/input"
	"github.com/jarmo/secrets-cli/v7/storage/path"
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
	printUsers(os.Stdout, users)
}

func printUsers(output io.Writer, users []crypto.VaultUser) {
	sort.Slice(users, func(i, j int) bool { return users[i].Name < users[j].Name })
	for _, user := range users {
		if user.PasswordBased {
			fmt.Fprintf(output, "%s (vault password)\n", user.Name)
		} else {
			fmt.Fprintln(output, user.Name)
		}
	}
}
