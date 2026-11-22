package path

import (
  "path/filepath"
  "github.com/pinzolo/xdgdir"
  "github.com/jarmo/secrets/storage/path"
  "github.com/jarmo/secrets/user"
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

func UsersDirectory() string {
  directory, err := path.UsersDirectory(configurationPath())
  if err != nil {
    panic(err)
  }
  if directory != "" {
    return directory
  }
  return user.DefaultDirectory(filepath.Dir(configurationPath()))
}

func Aliases() ([]string, error) {
  return path.Aliases(configurationPath())
}
