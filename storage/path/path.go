package path

import (
  "io/ioutil"
  "encoding/json"
  "github.com/pinzolo/xdgdir"
  "github.com/jarmo/secrets/storage/path"
)

func Get(alias string) (string, error) {
  if configs, err := path.Configurations(configurationPath()); err == nil {
    if confByAlias := path.FindByAlias(configs, alias); confByAlias != nil {
      return confByAlias.Path, nil
    } else {
      return configs[0].Path, nil
    }
  } else {
    return "", err
  }
}

func Store(vaultPath string, vaultAlias string) string {
  configurationPath := configurationPath()
  conf, _ := path.Configurations(configurationPath)
  conf = append(conf, path.Config{Path: vaultPath, Alias: vaultAlias})

  if configJSON, err := json.MarshalIndent(conf, "", " "); err != nil {
    panic(err)
  } else if err := ioutil.WriteFile(configurationPath, configJSON, 0600); err != nil {
    panic(err)
  }

  return configurationPath
}

func configurationPath() string {
  xdgApp := xdgdir.NewApp("secrets")
  xdgConfigurationFilePath, err := xdgApp.ConfigFile("config.json")
  if err != nil {
    panic(err)
  }

  return xdgConfigurationFilePath
}
