package cli

import (
  "github.com/docopt/docopt-go"
  "github.com/satori/go.uuid"
  "github.com/jarmo/secrets-cli/v6/cli/command"
)

func Command(version string, args []string) command.Executable {
  if len(args) >= 2 && args[0] == "__complete" {
    directory := ""
    if len(args) > 2 {
      directory = args[2]
    }
    return command.Completion{CandidateKind: args[1], UsersDirectory: directory}
  }
  arguments, _ := docopt.Parse(createUsage(), args, true, version, false)
  return createCommand(arguments)
}

func createUsage() string {
  return `secrets COMMAND [OPTIONS]

Usage:
  secrets (-h | --help)
  secrets (-v | --version)
  secrets list [FILTER] [--user=USERNAME] [--users-dir=USERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets add NAME [--user=USERNAME] [--users-dir=USERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets edit ID [--user=USERNAME] [--users-dir=USERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets delete ID [--user=USERNAME] [--users-dir=USERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets user-password-change USER [--users-dir=USERS_DIR]
  secrets initialize --path=VAULT_PATH --alias=VAULT_ALIAS --user=USERNAME [--users-dir=USERS_DIR]
  secrets user-create USER [--users-dir=USERS_DIR]
  secrets user-activate USER [--users-dir=USERS_DIR]
  secrets user-add USER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--users-dir=USERS_DIR] [--user=USERNAME]
  secrets user-list [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--users-dir=USERS_DIR] [--user=USERNAME]
  secrets user-remove USER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--users-dir=USERS_DIR] [--user=USERNAME]
  secrets user-destroy USER --user=USERNAME [--users-dir=USERS_DIR]
  secrets completion (bash | zsh)

Options:
  --alias VAULT_ALIAS    Optional vault alias.
  --path VAULT_PATH      Optional vault path. Defaults to the path in configuration.
  --users-dir USERS_DIR    Optional user storage directory.
  --user USERNAME    Authenticate using this user account.
  -h --help              Show this screen.
  -v --version           Show version.`
}

func createCommand(arguments map[string]interface {}) command.Executable {
  vaultAlias := vaultAlias(arguments)
  vaultPath := vaultPath(arguments)
  authenticatedUser, _ := arguments["--user"].(string)
  usersDirectory, _ := arguments["--users-dir"].(string)

  if arguments["list"].(bool) {
    if filter, hasValue := arguments["FILTER"].(string); !hasValue {
      return command.List{Filter: "", VaultAlias: vaultAlias, VaultPath: vaultPath, User: authenticatedUser, UsersDirectory: usersDirectory}
    } else {
      return command.List{Filter: filter, VaultAlias: vaultAlias, VaultPath: vaultPath, User: authenticatedUser, UsersDirectory: usersDirectory}
    }
  } else if arguments["add"].(bool) {
    return command.Add{Name: arguments["NAME"].(string), VaultAlias: vaultAlias, VaultPath: vaultPath, User: authenticatedUser, UsersDirectory: usersDirectory}
  } else if arguments["edit"].(bool) {
    id, _ := uuid.FromString(arguments["ID"].(string))
    return command.Edit{Id: id, VaultAlias: vaultAlias, VaultPath: vaultPath, User: authenticatedUser, UsersDirectory: usersDirectory}
  } else if arguments["delete"].(bool) {
    id, _ := uuid.FromString(arguments["ID"].(string))
    return command.Delete{Id: id, VaultAlias: vaultAlias, VaultPath: vaultPath, User: authenticatedUser, UsersDirectory: usersDirectory}
  } else if arguments["user-password-change"].(bool) {
    return command.UserPasswordChange{Username: arguments["USER"].(string), UsersDirectory: usersDirectory}
  } else if arguments["user-create"].(bool) {
    directory, _ := arguments["--users-dir"].(string)
    return command.UserCreate{Username: arguments["USER"].(string), UsersDirectory: directory}
  } else if arguments["user-activate"].(bool) {
    directory, _ := arguments["--users-dir"].(string)
    return command.UserActivate{Username: arguments["USER"].(string), UsersDirectory: directory}
  } else if arguments["user-add"].(bool) {
    directory, _ := arguments["--users-dir"].(string)
    return command.UserAdd{Username: arguments["USER"].(string), VaultAlias: vaultAlias, VaultPath: vaultPath, UsersDirectory: directory, User: authenticatedUser}
  } else if arguments["user-list"].(bool) {
    directory, _ := arguments["--users-dir"].(string)
    return command.UserList{VaultAlias: vaultAlias, VaultPath: vaultPath, UsersDirectory: directory, User: authenticatedUser}
  } else if arguments["user-remove"].(bool) {
    directory, _ := arguments["--users-dir"].(string)
    return command.UserRemove{Username: arguments["USER"].(string), VaultAlias: vaultAlias, VaultPath: vaultPath, UsersDirectory: directory, User: authenticatedUser}
  } else if arguments["user-destroy"].(bool) {
    return command.UserDestroy{Username: arguments["USER"].(string), User: authenticatedUser, UsersDirectory: usersDirectory}
  } else if arguments["completion"].(bool) {
    if arguments["bash"].(bool) {
      return command.Completion{Shell: "bash"}
    }
    return command.Completion{Shell: "zsh"}
  } else if arguments["initialize"].(bool) {
    return command.Initialize{VaultAlias: vaultAlias, VaultPath: vaultPath, User: authenticatedUser, UsersDirectory: usersDirectory}
  } else {
    return nil
  }
}

func vaultPath(arguments map[string]interface {}) string {
  if vaultPath, hasValue := arguments["--path"].(string); hasValue {
    return vaultPath
  } else {
    return ""
  }
}

func vaultAlias(arguments map[string]interface {}) string {
  if vaultAlias, hasValue := arguments["--alias"].(string); hasValue {
    return vaultAlias
  } else {
    return ""
  }
}
