package command

import (
  "fmt"
  "os"
  "github.com/jarmo/secrets/vault"
  "github.com/jarmo/secrets-cli/v6/cli/vaultfile"
  "github.com/jarmo/secrets-cli/v6/input"
)

type Add struct {
  Name string
  VaultPath string
  VaultAlias string
  User string
  UsersDirectory string
}

func (command Add) Execute() {
  openedVault := vaultfile.Open(vaultfile.Path(command.VaultAlias, command.VaultPath), command.User, command.UsersDirectory)
  secretName := command.Name
  secretValue := input.AskMultiline(fmt.Sprintf("Enter value for '%s':\n", command.Name))
  newSecret, newSecrets := vault.Add(openedVault.Secrets, secretName, secretValue)
  if err := openedVault.Write(newSecrets); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  fmt.Printf("Added: %s (%s)\n", newSecret.Name, newSecret.Id)
}

