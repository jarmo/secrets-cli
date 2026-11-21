package command

import (
	"crypto/ecdh"
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v7/cli/vaultfile"
	"github.com/jarmo/secrets-cli/v7/input"
	"github.com/jarmo/secrets-cli/v7/storage/path"
	"github.com/jarmo/secrets/user"
	"github.com/jarmo/secrets/storage"
)

type UserAdd struct {
	Username         string
	VaultAlias       string
	VaultPath        string
	UsersDirectory string
	User           string
}

func (command UserAdd) Execute() {
	directory := command.UsersDirectory
	if directory == "" {
		directory = path.UsersDirectory()
	}
	vaultPath := vaultfile.Path(command.VaultAlias, command.VaultPath)
	var err error
	if command.User != "" {
		openedVault := vaultfile.Open(vaultPath, command.User, directory)
		err = addUserForUser(vaultPath, directory, command.Username, openedVault.PrivateKey())
	} else {
		err = addUser(vaultPath, directory, command.Username, input.AskVaultPassword())
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("User %s added to vault!\n", command.Username)
}

func addUser(vaultPath, directory, username string, password []byte) error {
	userRecord, err := user.Load(directory, username)
	if err != nil {
		return err
	}
	return storage.AssociateUserRetiringPasswords(vaultPath, directory, password, userRecord.Identity)
}

func addUserForUser(vaultPath, directory, username string, privateKey *ecdh.PrivateKey) error {
	userRecord, err := user.Load(directory, username)
	if err != nil {
		return err
	}
	return storage.AssociateUserForUserRetiringPasswords(vaultPath, directory, privateKey, userRecord.Identity)
}
