package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v7/storage/path"
	"github.com/jarmo/secrets/user"
)

type UserCreate struct {
	Username         string
	UsersDirectory string
}

func (command UserCreate) Execute() {
	directory := command.UsersDirectory
	if directory == "" {
		directory = path.UsersDirectory()
	}
	code, expiresAt, err := user.CreatePending(directory, command.Username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Setup code for %s (expires %s):\n%s\n", command.Username, expiresAt.Format("2006-01-02 15:04 MST"), code)
}
