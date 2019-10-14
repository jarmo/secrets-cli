package command

import (
  "fmt"
  "os"
  "github.com/jarmo/secrets/vault"
  coreInput "github.com/jarmo/secrets/input"
  "github.com/jarmo/secrets-cli/v5/cli/vaultfile"
  "github.com/jarmo/secrets-cli/v5/input"
)

type ChangePassword struct {
  VaultPath string
  VaultAlias string
}

func (command ChangePassword) Execute() {
  currentPassword := input.AskVaultPassword()
  newPassword := coreInput.AskPassword("Enter new vault password: ")
  newPasswordConfirmation := coreInput.AskPassword("Enter new vault password again: ")

  if err := vault.ChangePassword(vaultfile.Path(command.VaultAlias, command.VaultPath), currentPassword, newPassword, newPasswordConfirmation); err != nil {
    fmt.Println(err)
    os.Exit(1)
  } else {
    fmt.Println("Vault password successfully changed!")
  }
}

