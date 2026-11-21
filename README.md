# secrets-cli

**Secure** and simple passwords manager written in [Go](https://golang.org/). It aims to be *NYAPM* (Not Yet Another Password Manager), but tries to be different from others by following UNIX philosophy of doing only one thing and doing it well.

This repository is for command-line client. There exists also a [self-hosted web solution](https://github.com/jarmo/secrets-web). Read more about [secrets](https://github.com/jarmo/secrets) in here.

## Installation

Download latest binary from [releases](https://github.com/jarmo/secrets-cli/releases), extract it and add it to somewhere in your **PATH**. That's it.

*Of course, you're free to compile your own version of binary to be 100% sure that it has not been tampered with, since this is an open-source project after all.*

## Usage

Here's an output from `secrets --help` command.

```
$ secrets COMMAND [OPTIONS]

Usage:
  secrets list [FILTER] [--user=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets add NAME [--user=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets edit ID [--user=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets delete ID [--user=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets user-password-change USER [--users-dir=USERS_DIR]
  secrets initialize --path=VAULT_PATH --alias=VAULT_ALIAS --user=USERNAME [--users-dir=USERS_DIR]
  secrets user-create USER [--users-dir=USERS_DIR]
  secrets user-activate USER [--users-dir=USERS_DIR]
  secrets user-add USER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--users-dir=USERS_DIR] [--user=USERNAME]
  secrets user-list [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--users-dir=USERS_DIR] [--user=USERNAME]
  secrets user-remove USER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--users-dir=USERS_DIR] [--user=USERNAME]
  secrets user-destroy USER --user=USERNAME [--users-dir=USERS_DIR]

Options:
  --alias VAULT_ALIAS    Vault alias.
  --path VAULT_PATH      Vault path.
  --users-dir USERS_DIR    User storage directory.
  --user USERNAME    Authenticate using this user account.
  -h --help              Show this screen.
  -v --version           Show version.
```

### Initializing Vault

After creating and activating a user account (see **User setup** below), create a vault and save its location and alias to a configuration file (supporting [XDG Base Directory standard](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)):

```
$ secrets initialize --path /home/user/.secrets.json --alias main --user alice
Enter user password: [Alice's password]
Vault successfully created and configured at /home/user/.config/secrets/config.json!
```

### Adding a New Secret

Add your first secret:

```
$ secrets add "my secret" --user alice
Enter user password: [Alice's password]
Enter value for 'my secret':
my secret value
Added: 
[299ed462-b171-4d67-ba21-264b221d9913]
my secret
my secret value
```

Because values can have multiple lines, you can enter whatever you want. Use ctrl+d on **macOS** and **Linux** or ctrl+z on **Windows** to complete entering multi-line values.

### Listing All Secrets

```
$ secrets list --user alice
Enter user password: [Alice's password]

[299ed462-b171-4d67-ba21-264b221d9913]
my secret
my secret value
```

### Listing Specific Secrets

```
$ secrets list "secret" --user alice
Enter user password: [Alice's password]

[299ed462-b171-4d67-ba21-264b221d9913]
my secret
my secret value
```

### Editing a Secret

```
$ secrets edit 299ed462-b171-4d67-ba21-264b221d9913 --user alice
Enter user password: [Alice's password]
Enter new name: different secret name
Enter new value:
different secret value
yet another secret value line
Edited: 
[299ed462-b171-4d67-ba21-264b221d9913]
different secret name
different secret value
yet another secret value line
```

### Deleting a Secret

```
$ secrets delete 299ed462-b171-4d67-ba21-264b221d9913 --user alice
Enter user password: 
Deleted: 
[299ed462-b171-4d67-ba21-264b221d9913]
different secret name
different secret value
yet another secret value line
```

### User setup

Create a user and share the generated setup code securely:

```
$ secrets user-create bob
Setup code for bob (expires ...):
[setup code]
```

Bob activates the user account by choosing a password:

```
$ secrets user-activate bob
Enter setup code: [setup code]
Enter new user password: [password]
Enter new user password again: [password]
User bob activated!
```

Setup codes expire after 48 hours. Both commands accept `--users-dir` to use the same user directory.

To change a user password without changing access to any vault:

```
$ secrets user-password-change bob
Enter current user password: [current password]
Enter new user password: [new password]
Enter new user password again: [new password]
User password successfully changed!
```

The same command supports `--users-dir`.

Once activated, a user can be granted access to a vault:

```
$ secrets user-add bob --alias production --user alice
Enter user password: [Alice's password]
User bob added to vault!
```

`user-add` also accepts `--path` and `--users-dir`.

Once granted access, a user can use their own credentials with `--user`:

```
$ secrets list --alias production --user bob
Enter user password: [Bob's password]
```

The same option works with `add`, `edit`, `delete`, `user-add`, `user-list` and `user-remove`. Use `--users-dir` if the accounts are stored outside the default directory.

List users or revoke a user's access to a vault:

```
$ secrets user-list --alias production --user alice
alice
bob

$ secrets user-remove bob --alias production --user alice
Enter user password: [Alice's password]
User bob removed from vault!
```

Removing a user rotates the vault encryption key while retaining access for other users. Both commands accept `--path` and `--users-dir`. The user's account is not deleted; it can be added to a vault again. Old backups can still be opened with keys that were valid at the time.

To revoke a user from all configured vaults and delete the account:

```
$ secrets user-destroy bob --user alice
Enter user password: [Alice's password]
User bob destroyed!
```

The acting user must have access to every configured vault. The command checks all vaults before making changes and does not delete the account unless all revocations succeed. All vaults containing the user must be in the CLI configuration. Backups retain historical access.

## Using multiple vaults

Just append `--alias` after any command to operate against selected vault.
When `--alias` is not specified a first vault existing in configuration file will be used.

## But how do I sync vault between different devices?!

One way to sync would be to use any already existing syncing platforms like Dropbox, Microsoft OneDrive or Google Drive.
Since you can specify vault storage location then it is up to you how (or if even) you sync.

## Legacy vaults

Starting with version 7.x, vaults use a new format that supports sharing between users. Legacy vaults are automatically converted to the new format when written. Support for the legacy format will be removed in the next major version.

To share a legacy vault, run `secrets user-add alice --alias production` and enter the original vault password. This grants Alice access and retires the vault password. For vaults already shared with users, running `user-add` for an existing user also retires any remaining vault-password access.

## Development

Retrieve dependencies, build and install binaries to `$GOPATH/bin/`

```
git clone https://github.com/jarmo/secrets-cli.git
cd secrets-cli
make
make install
```
