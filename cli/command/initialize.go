package command

import (
  "fmt"
  "os"
  "github.com/jarmo/secrets/member"
  "github.com/jarmo/secrets/storage"
  "github.com/jarmo/secrets-cli/v6/input"
  "github.com/jarmo/secrets-cli/v6/storage/path"
)

type Initialize struct {
  VaultPath string
  VaultAlias string
  Member string
  MembersDirectory string
}

func (command Initialize) Execute() {
  directory := command.MembersDirectory
  if directory == "" {
    directory = path.MembersDirectory()
  }
  password := input.AskPassword("Enter member password: ")
  defer func() {
    for index := range password {
      password[index] = 0
    }
  }()
  if err := initializeVault(command.VaultPath, directory, command.Member, password); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  configurationPath := path.Store(command.VaultPath, command.VaultAlias)
  fmt.Println(fmt.Sprintf("Vault successfully created and configured at %s!", configurationPath))
}

func initializeVault(vaultPath, membersDirectory, username string, password []byte) error {
  record, _, err := member.Authenticate(membersDirectory, username, password)
  if err != nil {
    return err
  }
  return storage.InitializeForMember(vaultPath, record.Identity)
}
