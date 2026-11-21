# secrets-cli

Command-line client for [secrets](https://github.com/jarmo/secrets). It works independently of the [web app](https://github.com/jarmo/secrets-web) and can also administer vaults used by the web app.

## Installation

Download a binary from [releases](https://github.com/jarmo/secrets-cli/releases), extract it, and put the `secrets` executable in your `PATH`. You can also build it from source.

## Getting started

Create a user first. The command prints a one-time setup code, valid for 48 hours:

```sh
secrets user-create alice
```

Activate the user with that code and choose a password (12–1,024 bytes):

```sh
secrets user-activate alice
```

Now create a vault for Alice:

```sh
secrets initialize --path /path/to/production.json --alias production --user alice
```

`initialize` creates an encrypted vault immediately. Alice's user password unlocks her identity, which grants access to the vault; there is no separate vault password.

Add and list secrets:

```sh
secrets add "my secret" --alias production --user alice
secrets list --alias production --user alice
secrets list "secret" --alias production --user alice
```

Use the secret ID returned by `list` to edit or delete it:

```sh
secrets edit SECRET_ID --alias production --user alice
secrets delete SECRET_ID --alias production --user alice
```

Secret values can span multiple lines. Finish input with Ctrl+D on Linux/macOS or Ctrl+Z on Windows.

## Users and vault sharing

Create and activate additional users in the same way. Once Bob has activated his account, Alice can grant him access without knowing his password:

```sh
secrets user-create bob
secrets user-activate bob
secrets user-add bob --alias production --user alice
```

Bob can then access the same vault with his own password:

```sh
secrets list --alias production --user bob
```

Manage access and account passwords:

```sh
secrets user-list --alias production --user alice
secrets user-remove bob --alias production --user alice
secrets user-password-change alice
```

`user-remove` revokes access to that vault but leaves Bob's account intact. To revoke a user from all configured vaults and delete the account:

```sh
secrets user-destroy bob --user alice
```

The acting user must be able to access the configured vaults. All vaults involving Bob must be present in the CLI configuration before destroying his account. Revocation rotates affected vault keys, but older backups remain decryptable with credentials valid when they were created.

## Vaults and configuration

`--alias` selects a configured vault; `--path` selects a vault file directly. Each vault has its own set of authorized users.

By default, the CLI stores vault aliases in `~/.config/secrets/config.json` and users in `~/.config/secrets/users`, respecting `XDG_CONFIG_HOME`. Use `--users-dir` to override the user directory for commands that access accounts.

Vault files can be backed up or synchronized using your preferred tools. Back up the user directory as well: without the users' encrypted identities and their passwords, encrypted vaults may become inaccessible.

Run `secrets --help` for the full command and option reference.

## Legacy vaults

Starting with version 7.x, vaults use a format that supports multiple users. Older vaults can still be read. Adding the first user with `secrets user-add alice --alias production` and the original vault password converts the vault to the new format and retires password-only access. A normal write also converts the legacy format.

A regular write also upgrades an old vault format, but the old password remains valid until user-based access is established with `user-add`. Back up old vaults before conversion. Support for the legacy format will be removed in the next major version.

## Development

```sh
git clone https://github.com/jarmo/secrets-cli.git
cd secrets-cli
make test
make install
```

`make install` builds the `secrets` executable into the Go binary directory. Use `INSTALL_DIR=/custom/bin` to change the installation directory.
