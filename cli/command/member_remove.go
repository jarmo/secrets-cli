package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/cli/vaultfile"
	"github.com/jarmo/secrets-cli/v6/input"
	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/storage"
)

type MemberRemove struct {
	Username         string
	VaultAlias       string
	VaultPath        string
	MembersDirectory string
	Member           string
}

func (command MemberRemove) Execute() {
	directory := command.MembersDirectory
	if directory == "" {
		directory = path.MembersDirectory()
	}
	vaultPath := vaultfile.Path(command.VaultAlias, command.VaultPath)
	var err error
	if command.Member != "" {
		openedVault := vaultfile.Open(vaultPath, command.Member, directory)
		err = storage.RemoveMemberForMember(vaultPath, directory, command.Username, openedVault.PrivateKey())
	} else {
		err = storage.RemoveMember(vaultPath, directory, command.Username, input.AskVaultPassword())
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Member %s removed from vault!\n", command.Username)
}
