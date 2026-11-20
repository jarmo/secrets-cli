package cli

import (
  "testing"
  "github.com/satori/go.uuid"
  "github.com/jarmo/secrets-cli/v6/cli/command"
)

const version = "1.3.3.7"

func TestExecute_ListWithoutFilter(t *testing.T) {
  filter := ""

  switch parsedCommand := Command(version, []string{"list", filter}).(type) {
    case command.List:
      if parsedCommand.Filter != filter {
        t.Fatalf("Expected filter to be '%v' but was: '%v'", filter, parsedCommand.Filter)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ListWithoutFilterAndWithCustomVaultPath(t *testing.T) {
  filter := ""
  vaultPath := "/foo/bar/baz"

  switch parsedCommand := Command(version, []string{"list", "--path", vaultPath}).(type) {
    case command.List:
      if parsedCommand.Filter != filter {
        t.Fatalf("Expected filter to be '%v' but was: '%v'", filter, parsedCommand.Filter)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != vaultPath {
        t.Fatalf("Expected VaultPath to be '%v' but was: '%v'", vaultPath, parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ListWithoutFilterAndWithVaultAlias(t *testing.T) {
  filter := ""
  vaultAlias := "foo-bar"

  switch parsedCommand := Command(version, []string{"list", "--alias", vaultAlias}).(type) {
    case command.List:
      if parsedCommand.Filter != filter {
        t.Fatalf("Expected filter to be '%v' but was: '%v'", filter, parsedCommand.Filter)
      }
      if parsedCommand.VaultAlias != vaultAlias {
        t.Fatalf("Expected VaultAlias to be '%v' but was: '%v'", vaultAlias, parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ListWithFilter(t *testing.T) {
  filter := "custom-filter"

  switch parsedCommand := Command(version, []string{"list", filter}).(type) {
    case command.List:
      if parsedCommand.Filter != filter {
        t.Fatalf("Expected filter to be '%v', but was '%v'", filter, parsedCommand.Filter)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ListWithFilterAndWithCustomVaultPath(t *testing.T) {
  filter := "custom-filter"
  vaultPath := "/foo/bar/baz"

  switch parsedCommand := Command(version, []string{"list", filter, "--path", vaultPath}).(type) {
    case command.List:
      if parsedCommand.Filter != filter {
        t.Fatalf("Expected filter to be '%v', but was '%v'", filter, parsedCommand.Filter)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != vaultPath {
        t.Fatalf("Expected VaultPath to be '%v' but was: '%v'", vaultPath, parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ListWithFilterAndWithVaultAlias(t *testing.T) {
  filter := "custom-filter"
  vaultAlias := "foo-bar"

  switch parsedCommand := Command(version, []string{"list", filter, "--alias", vaultAlias}).(type) {
    case command.List:
      if parsedCommand.Filter != filter {
        t.Fatalf("Expected filter to be '%v', but was '%v'", filter, parsedCommand.Filter)
      }
      if parsedCommand.VaultAlias != vaultAlias {
        t.Fatalf("Expected VaultAlias to be '%v' but was: '%v'", vaultAlias, parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_Add(t *testing.T) {
  name := "custom-name"

  switch parsedCommand := Command(version, []string{"add", name}).(type) {
    case command.Add:
      if parsedCommand.Name != name {
        t.Fatalf("Expected name to be '%v', but was '%v'", name, parsedCommand.Name)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_AddWithCustomVaultPath(t *testing.T) {
  name := "custom-name"
  vaultPath := "/foo/bar/baz"

  switch parsedCommand := Command(version, []string{"add", name, "--path", vaultPath}).(type) {
    case command.Add:
      if parsedCommand.Name != name {
        t.Fatalf("Expected name to be '%v', but was '%v'", name, parsedCommand.Name)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != vaultPath {
        t.Fatalf("Expected VaultPath to be '%v' but was: '%v'", vaultPath, parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_AddWithVaultAlias(t *testing.T) {
  name := "custom-name"
  vaultAlias := "foo-bar"

  switch parsedCommand := Command(version, []string{"add", name, "--alias", vaultAlias}).(type) {
    case command.Add:
      if parsedCommand.Name != name {
        t.Fatalf("Expected name to be '%v', but was '%v'", name, parsedCommand.Name)
      }
      if parsedCommand.VaultAlias != vaultAlias {
        t.Fatalf("Expected VaultAlias to be '%v' but was: '%v'", vaultAlias, parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_Edit(t *testing.T) {
  id := uuid.NewV4()

  switch parsedCommand := Command(version, []string{"edit", id.String()}).(type) {
    case command.Edit:
      if parsedCommand.Id != id {
        t.Fatalf("Expected id to be '%v', but was '%v'", id, parsedCommand.Id)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_EditWithCustomVaultPath(t *testing.T) {
  id := uuid.NewV4()
  vaultPath := "/foo/bar/baz"

  switch parsedCommand := Command(version, []string{"edit", id.String(), "--path", vaultPath}).(type) {
    case command.Edit:
      if parsedCommand.Id != id {
        t.Fatalf("Expected id to be '%v', but was '%v'", id, parsedCommand.Id)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != vaultPath {
        t.Fatalf("Expected VaultPath to be '%v' but was: '%v'", vaultPath, parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_EditWithVaultAlias(t *testing.T) {
  id := uuid.NewV4()
  vaultAlias := "foo-bar"

  switch parsedCommand := Command(version, []string{"edit", id.String(), "--alias", vaultAlias}).(type) {
    case command.Edit:
      if parsedCommand.Id != id {
        t.Fatalf("Expected id to be '%v', but was '%v'", id, parsedCommand.Id)
      }
      if parsedCommand.VaultAlias != vaultAlias {
        t.Fatalf("Expected VaultAlias to be '%v' but was: '%v'", vaultAlias, parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_Delete(t *testing.T) {
  id := uuid.NewV4()

  switch parsedCommand := Command(version, []string{"delete", id.String()}).(type) {
    case command.Delete:
      if parsedCommand.Id != id {
        t.Fatalf("Expected id to be '%v', but was '%v'", id, parsedCommand.Id)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_DeleteWithCustomVaultPath(t *testing.T) {
  id := uuid.NewV4()
  vaultPath := "/foo/bar/baz"

  switch parsedCommand := Command(version, []string{"delete", id.String(), "--path", vaultPath}).(type) {
    case command.Delete:
      if parsedCommand.Id != id {
        t.Fatalf("Expected id to be '%v', but was '%v'", id, parsedCommand.Id)
      }
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != vaultPath {
        t.Fatalf("Expected VaultPath to be '%v' but was: '%v'", vaultPath, parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_DeleteWithVaultAlias(t *testing.T) {
  id := uuid.NewV4()
  vaultAlias := "foo-bar"

  switch parsedCommand := Command(version, []string{"delete", id.String(), "--alias", vaultAlias}).(type) {
    case command.Delete:
      if parsedCommand.Id != id {
        t.Fatalf("Expected id to be '%v', but was '%v'", id, parsedCommand.Id)
      }
      if parsedCommand.VaultAlias != vaultAlias {
        t.Fatalf("Expected VaultAlias to be '%v' but was: '%v'", vaultAlias, parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ChangePassword(t *testing.T) {
  switch parsedCommand := Command(version, []string{"change-password"}).(type) {
    case command.ChangePassword:
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ChangePasswordWithCustomVaultPath(t *testing.T) {
  vaultPath := "/foo/bar/baz"

  switch parsedCommand := Command(version, []string{"change-password", "--path", vaultPath}).(type) {
    case command.ChangePassword:
      if parsedCommand.VaultAlias != "" {
        t.Fatalf("Expected VaultAlias to be empty but was: '%v'", parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != vaultPath {
        t.Fatalf("Expected VaultPath to be '%v' but was: '%v'", vaultPath, parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_ChangePasswordWithVaultAlias(t *testing.T) {
  vaultAlias := "foo-bar"

  switch parsedCommand := Command(version, []string{"change-password", "--alias", vaultAlias}).(type) {
    case command.ChangePassword:
      if parsedCommand.VaultAlias != vaultAlias {
        t.Fatalf("Expected VaultAlias to be '%v' but was: '%v'", vaultAlias, parsedCommand.VaultAlias)
      }
      if parsedCommand.VaultPath != "" {
        t.Fatalf("Expected VaultPath to be empty but was: '%v'", parsedCommand.VaultPath)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_Initialize(t *testing.T) {
  vaultPath := "/foo/bar/baz"
  vaultAlias := "foo-bar"

  switch parsedCommand := Command(version, []string{"initialize", "--path", vaultPath, "--alias", vaultAlias, "--user", "alice", "--users-dir", "/tmp/users"}).(type) {
    case command.Initialize:
      if parsedCommand.VaultPath != vaultPath {
        t.Fatalf("Expected vault path to be '%v', but was '%v'", vaultPath, parsedCommand.VaultPath)
      }
      if parsedCommand.VaultAlias != vaultAlias {
        t.Fatalf("Expected VaultAlias to be '%v' but was: '%v'", vaultAlias, parsedCommand.VaultAlias)
      }
      if parsedCommand.User != "alice" || parsedCommand.UsersDirectory != "/tmp/users" {
        t.Fatalf("Expected user authentication during initialization: %+v", parsedCommand)
      }
    default:
      t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}

func TestExecute_UserCreate(t *testing.T) {
  parsedCommand, ok := Command(version, []string{"user-create", "bob"}).(command.UserCreate)
  if !ok || parsedCommand.Username != "bob" || parsedCommand.UsersDirectory != "" {
    t.Fatalf("Unexpected user-create command: %#v", parsedCommand)
  }

  parsedCommand, ok = Command(version, []string{"user-create", "bob", "--users-dir", "/var/lib/secrets/users"}).(command.UserCreate)
  if !ok || parsedCommand.Username != "bob" || parsedCommand.UsersDirectory != "/var/lib/secrets/users" {
    t.Fatalf("Unexpected user-create options: %#v", parsedCommand)
  }
}

func TestExecute_UserActivate(t *testing.T) {
  parsedCommand, ok := Command(version, []string{"user-activate", "bob"}).(command.UserActivate)
  if !ok || parsedCommand.Username != "bob" || parsedCommand.UsersDirectory != "" {
    t.Fatalf("Unexpected user-activate command: %#v", parsedCommand)
  }

  parsedCommand, ok = Command(version, []string{"user-activate", "bob", "--users-dir", "/var/lib/secrets/users"}).(command.UserActivate)
  if !ok || parsedCommand.Username != "bob" || parsedCommand.UsersDirectory != "/var/lib/secrets/users" {
    t.Fatalf("Unexpected user-activate options: %#v", parsedCommand)
  }
}

func TestExecute_UserAdd(t *testing.T) {
  parsedCommand, ok := Command(version, []string{"user-add", "bob"}).(command.UserAdd)
  if !ok || parsedCommand.Username != "bob" || parsedCommand.VaultAlias != "" || parsedCommand.VaultPath != "" || parsedCommand.UsersDirectory != "" {
    t.Fatalf("Unexpected user-add command: %#v", parsedCommand)
  }

  parsedCommand, ok = Command(version, []string{"user-add", "bob", "--alias", "production", "--users-dir", "/private/users"}).(command.UserAdd)
  if !ok || parsedCommand.Username != "bob" || parsedCommand.VaultAlias != "production" || parsedCommand.VaultPath != "" || parsedCommand.UsersDirectory != "/private/users" {
    t.Fatalf("Unexpected user-add options: %#v", parsedCommand)
  }

  parsedCommand, ok = Command(version, []string{"user-add", "bob", "--path", "/private/vault.json"}).(command.UserAdd)
  if !ok || parsedCommand.VaultPath != "/private/vault.json" || parsedCommand.VaultAlias != "" {
    t.Fatalf("Unexpected user-add vault path: %#v", parsedCommand)
  }
}

func TestExecute_UserListAndRemove(t *testing.T) {
  list := Command(version, []string{"user-list", "--alias", "production", "--users-dir", "/tmp/users"})
  listed, ok := list.(command.UserList)
  if !ok || listed.VaultAlias != "production" || listed.UsersDirectory != "/tmp/users" {
    t.Fatalf("Unexpected user-list command: %+v", list)
  }
  remove := Command(version, []string{"user-remove", "bob", "--path", "/tmp/vault", "--users-dir", "/tmp/users"})
  removed, ok := remove.(command.UserRemove)
  if !ok || removed.Username != "bob" || removed.VaultPath != "/tmp/vault" || removed.UsersDirectory != "/tmp/users" {
    t.Fatalf("Unexpected user-remove command: %+v", remove)
  }
}

func TestParseUserAuthenticationForVaultCommands(t *testing.T) {
	for _, args := range [][]string{
		{"list", "--alias=production", "--user=alice", "--users-dir=/tmp/users"},
		{"add", "example", "--user=alice", "--users-dir=/tmp/users"},
		{"edit", uuid.NewV4().String(), "--user=alice", "--users-dir=/tmp/users"},
		{"delete", uuid.NewV4().String(), "--user=alice", "--users-dir=/tmp/users"},
		{"user-add", "bob", "--user=alice", "--users-dir=/tmp/users"},
		{"user-list", "--user=alice", "--users-dir=/tmp/users"},
		{"user-remove", "bob", "--user=alice", "--users-dir=/tmp/users"},
	} {
		parsed := Command(version, args)
		var userName, usersDirectory string
		switch cmd := parsed.(type) {
		case command.List:
			userName, usersDirectory = cmd.User, cmd.UsersDirectory
		case command.Add:
			userName, usersDirectory = cmd.User, cmd.UsersDirectory
		case command.Edit:
			userName, usersDirectory = cmd.User, cmd.UsersDirectory
		case command.Delete:
			userName, usersDirectory = cmd.User, cmd.UsersDirectory
		case command.UserAdd:
			userName, usersDirectory = cmd.User, cmd.UsersDirectory
		case command.UserList:
			userName, usersDirectory = cmd.User, cmd.UsersDirectory
		case command.UserRemove:
			userName, usersDirectory = cmd.User, cmd.UsersDirectory
		default:
			t.Fatalf("Unexpected command type %T for %v", parsed, args)
		}
		if userName != "alice" || usersDirectory != "/tmp/users" {
			t.Fatalf("Invalid user options for %v: %q %q", args, userName, usersDirectory)
		}
	}
}

func TestExecute_UserDestroy(t *testing.T) {
  switch parsedCommand := Command(version, []string{"user-destroy", "bob", "--user", "alice"}).(type) {
  case command.UserDestroy:
    if parsedCommand.Username != "bob" || parsedCommand.User != "alice" {
      t.Fatalf("Unexpected user-destroy arguments: %+v", parsedCommand)
    }
  default:
    t.Fatalf("Got unexpected command: %T", parsedCommand)
  }
}
