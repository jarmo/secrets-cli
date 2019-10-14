package path

import (
  "github.com/pinzolo/xdgdir"
  "github.com/jarmo/secrets/storage/path"
)

func Get(alias string) (string, error) {
  if vaultPath, err := path.Get(configurationPath(), alias); err != nil {
    return "", err
  } else {
    return vaultPath, nil
  }
}

func Store(vaultPath string, vaultAlias string) string {
  return path.Store(configurationPath(), vaultPath, vaultAlias)
}

func configurationPath() string {
  xdgApp := xdgdir.NewApp("secrets")
  xdgConfigurationFilePath, err := xdgApp.ConfigFile("config.json")
  if err != nil {
    panic(err)
  }

  return xdgConfigurationFilePath
}
