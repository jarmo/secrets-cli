package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/cli/vaultfile"
	"github.com/jarmo/secrets-cli/v6/input"
	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/storage"
)

type MemberList struct {
	VaultAlias       string
	VaultPath        string
	MembersDirectory string
}

func (command MemberList) Execute() {
	directory := command.MembersDirectory
	if directory == "" {
		directory = path.MembersDirectory()
	}
	vaultPath := vaultfile.Path(command.VaultAlias, command.VaultPath)
	members, err := storage.ListMembers(vaultPath, directory, input.AskVaultPassword())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, member := range members {
		if member.PasswordBased {
			fmt.Printf("%s (vault password)\n", member.Name)
		} else {
			fmt.Println(member.Name)
		}
	}
}
