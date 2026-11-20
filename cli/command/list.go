package command

import (
  "fmt"
  "github.com/jarmo/secrets/vault"
  "github.com/jarmo/secrets-cli/v6/cli/vaultfile"
)

type List struct {
  Filter string
  VaultPath string
  VaultAlias string
  User string
  UsersDirectory string
}

func (command List) Execute() {
  openedVault := vaultfile.Open(vaultfile.Path(command.VaultAlias, command.VaultPath), command.User, command.UsersDirectory)
  for _, secret := range vault.List(openedVault.Secrets, command.Filter) {
    fmt.Println(secret)
  }
}

