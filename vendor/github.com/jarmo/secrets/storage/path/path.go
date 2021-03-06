package path

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type config struct {
	Path  string
	Alias string
}

func Get(configurationPath, alias string) (string, error) {
	if configs, err := configurations(configurationPath); err == nil {
		if confByAlias := findByAlias(configs, alias); confByAlias != nil {
			return confByAlias.Path, nil
		} else {
			return "", errors.New("Vault not found!")
		}
	} else {
		return "", err
	}
}

func Store(configurationPath, vaultPath, vaultAlias string) string {
	conf, _ := configurations(configurationPath)
	conf = append(conf, config{Path: vaultPath, Alias: vaultAlias})

	if configJSON, err := json.MarshalIndent(conf, "", " "); err != nil {
		panic(err)
	} else if err := ioutil.WriteFile(configurationPath, configJSON, 0600); err != nil {
		panic(err)
	}

	return configurationPath
}

func configurations(path string) ([]config, error) {
	if configJSON, err := ioutil.ReadFile(path); os.IsNotExist(err) {
		return make([]config, 0), errors.New("Vault is not configured!")
	} else {
		var conf []config
		if err := json.Unmarshal(configJSON, &conf); err == nil {
			return conf, nil
		} else {
			return make([]config, 0), err
		}
	}
}

func findByAlias(configs []config, alias string) *config {
	for _, config := range configs {
		if strings.ToLower(config.Alias) == strings.ToLower(alias) {
			return &config
		}
	}

	return nil
}
