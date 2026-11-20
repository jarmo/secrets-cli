package command

import (
  "fmt"
  "os"
  "github.com/satori/go.uuid"
  "github.com/jarmo/secrets/vault"
  "github.com/jarmo/secrets-cli/v6/cli/vaultfile"
)

type Delete struct {
  Id uuid.UUID
  VaultPath string
  VaultAlias string
  Member string
  MembersDirectory string
}

func (command Delete) Execute() {
  openedVault := vaultfile.Open(vaultfile.Path(command.VaultAlias, command.VaultPath), command.Member, command.MembersDirectory)
  deletedSecret, newSecrets, err := vault.Delete(openedVault.Secrets, command.Id)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  } else {
    if err := openedVault.Write(newSecrets); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
    fmt.Println("Deleted:", deletedSecret)
  }
}

