package command

import (
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v7/completion"
)

type Completion struct {
	Shell          string
	CandidateKind  string
	UsersDirectory string
}

func (command Completion) Execute() {
	if command.CandidateKind != "" {
		completion.Candidates(os.Stdout, command.CandidateKind, command.UsersDirectory)
		return
	}
	if err := completion.Script(os.Stdout, command.Shell); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
