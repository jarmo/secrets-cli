package command

import (
  "fmt"
  "os"
  "github.com/jarmo/secrets/user"
  "github.com/jarmo/secrets/storage"
  "github.com/jarmo/secrets-cli/v7/input"
  "github.com/jarmo/secrets-cli/v7/storage/path"
)

type Initialize struct {
  VaultPath string
  VaultAlias string
  User string
  UsersDirectory string
}

func (command Initialize) Execute() {
  directory := command.UsersDirectory
  if directory == "" {
    directory = path.UsersDirectory()
  }
  password := input.AskPassword("Enter user password: ")
  defer func() {
    for index := range password {
      password[index] = 0
    }
  }()
  if err := initializeVault(command.VaultPath, directory, command.User, password); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  configurationPath := path.Store(command.VaultPath, command.VaultAlias)
  fmt.Println(fmt.Sprintf("Vault successfully created and configured at %s!", configurationPath))
}

func initializeVault(vaultPath, usersDirectory, username string, password []byte) error {
  record, _, err := user.Authenticate(usersDirectory, username, password)
  if err != nil {
    return err
  }
  return storage.InitializeForUser(vaultPath, record.Identity)
}
