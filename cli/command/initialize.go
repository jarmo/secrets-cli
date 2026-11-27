package command

import (
  "fmt"
  "os"
  "github.com/jarmo/secrets/user"
  "github.com/jarmo/secrets/storage"
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
  if err := initializeVault(command.VaultPath, directory, command.User); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  configurationPath := path.Store(command.VaultPath, command.VaultAlias)
  fmt.Println(fmt.Sprintf("Vault successfully created and configured at %s!", configurationPath))
}

func initializeVault(vaultPath, usersDirectory, username string) error {
  record, err := user.Load(usersDirectory, username)
  if err != nil {
    return err
  }
  return storage.InitializeForUser(vaultPath, record.Identity)
}
