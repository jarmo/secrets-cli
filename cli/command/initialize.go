package command

import (
  "bytes"
  "errors"
  "fmt"
  "os"
  "github.com/jarmo/secrets/storage"
  "github.com/jarmo/secrets-cli/v6/input"
  "github.com/jarmo/secrets-cli/v6/storage/path"
)

type Initialize struct {
  VaultPath string
  VaultAlias string
}

func (command Initialize) Execute() {
  password := input.AskPassword("Enter new vault password: ")
  confirmation := input.AskPassword("Enter new vault password again: ")
  if err := initializeVault(command.VaultPath, password, confirmation); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  configurationPath := path.Store(command.VaultPath, command.VaultAlias)
  fmt.Println(fmt.Sprintf("Vault successfully created and configured at %s!", configurationPath))
}

func initializeVault(vaultPath string, password, confirmation []byte) error {
  if !bytes.Equal(password, confirmation) {
    return errors.New("Passwords do not match!")
  }
  return storage.Initialize(vaultPath, password)
}

