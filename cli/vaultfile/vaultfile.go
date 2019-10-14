package vaultfile

import (
  "fmt"
  "os"
  "github.com/jarmo/secrets/storage"
  "github.com/jarmo/secrets/secret"
  "github.com/jarmo/secrets-cli/v5/storage/path"
  "github.com/jarmo/secrets-cli/v5/input"
)

func Read(vaultPath string) ([]secret.Secret, string, []byte) {
  password := input.AskVaultPassword()
  secrets, err := storage.Read(vaultPath, password)

  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  return secrets, vaultPath, password
}

func Path(alias string, customPath string) string {
  if customPath != "" {
    return customPath
  } else {
    vaultPath, err := path.Get(alias)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    return vaultPath
  }
}

