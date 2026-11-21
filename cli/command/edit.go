package command

import (
  "fmt"
  "os"
  "github.com/satori/go.uuid"
  "github.com/jarmo/secrets/vault"
  "github.com/jarmo/secrets-cli/v6/input"
  "github.com/jarmo/secrets-cli/v6/cli/vaultfile"
)

type Edit struct {
  Id uuid.UUID
  VaultPath string
  VaultAlias string
  User string
  UsersDirectory string
}

func (command Edit) Execute() {
  openedVault := vaultfile.Open(vaultfile.Path(command.VaultAlias, command.VaultPath), command.User, command.UsersDirectory)
  newName := input.Ask(fmt.Sprintf("Enter new name: "))
  newValue := input.AskMultiline("Enter new value:\n")
  editedSecret, newSecrets, err := vault.Edit(openedVault.Secrets, command.Id, newName, newValue)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  } else {
    if err := openedVault.Write(newSecrets); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
    fmt.Printf("Edited: %s (%s)\n", editedSecret.Name, editedSecret.Id)
  }
}

