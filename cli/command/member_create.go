package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/member"
)

type MemberCreate struct {
	Username         string
	MembersDirectory string
}

func (command MemberCreate) Execute() {
	directory := command.MembersDirectory
	if directory == "" {
		directory = path.MembersDirectory()
	}
	code, expiresAt, err := member.CreatePending(directory, command.Username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Setup code for %s (expires %s):\n%s\n", command.Username, expiresAt.Format("2006-01-02 15:04 MST"), code)
}
