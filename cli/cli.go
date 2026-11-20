package cli

import (
  "github.com/docopt/docopt-go"
  "github.com/satori/go.uuid"
  "github.com/jarmo/secrets-cli/v6/cli/command"
)

func Command(version string, args []string) command.Executable {
  arguments, _ := docopt.Parse(createUsage(), args, true, version, false)
  return createCommand(arguments)
}

func createUsage() string {
  return `secrets COMMAND [OPTIONS]

Usage:
  secrets list [FILTER] [--member=USERNAME] [--members-dir=MEMBERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets add NAME [--member=USERNAME] [--members-dir=MEMBERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets edit ID [--member=USERNAME] [--members-dir=MEMBERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets delete ID [--member=USERNAME] [--members-dir=MEMBERS_DIR] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets change-password [--alias=VAULT_PATH | --path=VAULT_PATH]
  secrets initialize --path=VAULT_PATH --alias=VAULT_ALIAS
  secrets member-create MEMBER [--members-dir=MEMBERS_DIR]
  secrets member-activate MEMBER [--members-dir=MEMBERS_DIR]
  secrets member-add MEMBER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--members-dir=MEMBERS_DIR] [--member=USERNAME]
  secrets member-list [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--members-dir=MEMBERS_DIR] [--member=USERNAME]
  secrets member-remove MEMBER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--members-dir=MEMBERS_DIR] [--member=USERNAME]

Options:
  --alias VAULT_ALIAS    Optional vault alias.
  --path VAULT_PATH      Optional vault path. Defaults to the path in configuration.
  --members-dir MEMBERS_DIR    Optional member storage directory.
  --member USERNAME    Authenticate using this member account.
  -h --help              Show this screen.
  -v --version           Show version.`
}

func createCommand(arguments map[string]interface {}) command.Executable {
  vaultAlias := vaultAlias(arguments)
  vaultPath := vaultPath(arguments)
  authenticatedMember, _ := arguments["--member"].(string)
  membersDirectory, _ := arguments["--members-dir"].(string)

  if arguments["list"].(bool) {
    if filter, hasValue := arguments["FILTER"].(string); !hasValue {
      return command.List{Filter: "", VaultAlias: vaultAlias, VaultPath: vaultPath, Member: authenticatedMember, MembersDirectory: membersDirectory}
    } else {
      return command.List{Filter: filter, VaultAlias: vaultAlias, VaultPath: vaultPath, Member: authenticatedMember, MembersDirectory: membersDirectory}
    }
  } else if arguments["add"].(bool) {
    return command.Add{Name: arguments["NAME"].(string), VaultAlias: vaultAlias, VaultPath: vaultPath, Member: authenticatedMember, MembersDirectory: membersDirectory}
  } else if arguments["edit"].(bool) {
    id, _ := uuid.FromString(arguments["ID"].(string))
    return command.Edit{Id: id, VaultAlias: vaultAlias, VaultPath: vaultPath, Member: authenticatedMember, MembersDirectory: membersDirectory}
  } else if arguments["delete"].(bool) {
    id, _ := uuid.FromString(arguments["ID"].(string))
    return command.Delete{Id: id, VaultAlias: vaultAlias, VaultPath: vaultPath, Member: authenticatedMember, MembersDirectory: membersDirectory}
  } else if arguments["change-password"].(bool) {
    return command.ChangePassword{VaultAlias: vaultAlias, VaultPath: vaultPath}
  } else if arguments["member-create"].(bool) {
    directory, _ := arguments["--members-dir"].(string)
    return command.MemberCreate{Username: arguments["MEMBER"].(string), MembersDirectory: directory}
  } else if arguments["member-activate"].(bool) {
    directory, _ := arguments["--members-dir"].(string)
    return command.MemberActivate{Username: arguments["MEMBER"].(string), MembersDirectory: directory}
  } else if arguments["member-add"].(bool) {
    directory, _ := arguments["--members-dir"].(string)
    return command.MemberAdd{Username: arguments["MEMBER"].(string), VaultAlias: vaultAlias, VaultPath: vaultPath, MembersDirectory: directory, Member: authenticatedMember}
  } else if arguments["member-list"].(bool) {
    directory, _ := arguments["--members-dir"].(string)
    return command.MemberList{VaultAlias: vaultAlias, VaultPath: vaultPath, MembersDirectory: directory, Member: authenticatedMember}
  } else if arguments["member-remove"].(bool) {
    directory, _ := arguments["--members-dir"].(string)
    return command.MemberRemove{Username: arguments["MEMBER"].(string), VaultAlias: vaultAlias, VaultPath: vaultPath, MembersDirectory: directory, Member: authenticatedMember}
  } else if arguments["initialize"].(bool) {
    return command.Initialize{VaultAlias: vaultAlias, VaultPath: vaultPath}
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
