package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/cli/vaultfile"
	"github.com/jarmo/secrets-cli/v6/input"
	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/member"
	"github.com/jarmo/secrets/storage"
)

type MemberAdd struct {
	Username         string
	VaultAlias       string
	VaultPath        string
	MembersDirectory string
}

func (command MemberAdd) Execute() {
	directory := command.MembersDirectory
	if directory == "" {
		directory = path.MembersDirectory()
	}
	vaultPath := vaultfile.Path(command.VaultAlias, command.VaultPath)
	password := input.AskVaultPassword()
	if err := addMember(vaultPath, directory, command.Username, password); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Member %s added to vault!\n", command.Username)
}

func addMember(vaultPath, directory, username string, password []byte) error {
	memberRecord, err := member.Load(directory, username)
	if err != nil {
		return err
	}
	return storage.AssociateMember(vaultPath, password, memberRecord.Identity)
}
