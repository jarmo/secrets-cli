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

By default, the CLI stores vault aliases in `~/.config/secrets/config.json` and users in `~/.config/secrets/users`, respecting `XDG_CONFIG_HOME`. To keep user accounts somewhere else (for example, in a synchronized Dropbox folder), set one global `UsersDir` in `config.json`:

```json
{
  "UsersDir": "/home/alice/Dropbox/secrets/users",
  "Vaults": [
    {"Path": "/home/alice/Dropbox/secrets/production.json", "Alias": "production"}
  ]
}
```

Use an absolute `UsersDir` path. The old JSON array of vault aliases remains supported and uses the default account directory. The CLI upgrades that format on its next configuration write; it does not rewrite configuration files merely to read them. The `--users-dir` flag overrides the configured directory for an individual command. Changing `UsersDir` does not move existing account files: copy or sync `users/` to the chosen location before using it. Update web to a version that understands the new format before changing a shared configuration file.

Vault files can be backed up or synchronized using your preferred tools. Back up the user directory as well: without the users' encrypted identities and their passwords, encrypted vaults may become inaccessible.

Run `secrets --help` for the full command and option reference. Check the installed version with `secrets --version`. Release builds use the Git tag; development builds show a commit identifier (and `-dirty` for uncommitted changes).

## Migrating legacy vaults

Version 7.x can convert older password-only vaults to user-based access. **Migrate every legacy vault before upgrading to the next major version**, which will remove legacy vault support. Adding the first user with `secrets user-add` converts that vault and retires its original vault password. A normal write may update the storage format without retiring password-based access, so it is not a substitute for `user-add`.

Before migrating, make secure backups of **all vault files**, `~/.config/secrets/config.json`, and `~/.config/secrets/users/` (or their `XDG_CONFIG_HOME` equivalents). Keep the old vault passwords needed to restore pre-migration backups.

```sh
# Skip these if alice already has an activated account.
secrets user-create alice
secrets user-activate alice

# Replace these with ALL legacy vault aliases in your config.json.
for alias in production personal staging; do
  secrets user-add alice --alias "$alias" || break
  secrets user-list --alias "$alias" --user alice || break
done
```

Each `user-add` prompts for that vault's **old vault password**; each `user-list` prompts for Alice's **user password** and verifies that the converted vault can be unlocked. Include only legacy vaults in the loop. For vaults not configured with an alias, repeat the commands using `--path /path/to/vault.json` instead. If any conversion or verification fails, resolve it before upgrading. For additional users, activate their accounts and grant access with `secrets user-add bob --alias production --user alice` (using Alice’s account password).

**Back up the migrated vault files and user-account data again after verification**, along with the configuration. Keep the pre-migration backups until you are confident recovery works, and protect both backup sets: the original backups still use the legacy format and will not be readable by the next major version.

## Shell completion

Generate completion scripts for Bash or Zsh. Completion suggests commands, options, configured vault aliases, and local usernames without unlocking vaults.

**Bash:**

```sh
mkdir -p ~/.local/share/bash-completion/completions
secrets completion bash > ~/.local/share/bash-completion/completions/secrets
```

Create the destination directory if needed. Bash completion needs to be enabled in your shell; start a new shell after installing the script.

**Zsh:**

```sh
mkdir -p ~/.zsh/completions
secrets completion zsh > ~/.zsh/completions/_secrets
```

Add the directory to `fpath` before initializing completion in `~/.zshrc`:

```sh
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit
compinit
```

Then start a new shell. Secret IDs are not completed because looking them up would require decrypting a vault.

## Development

```sh
git clone https://github.com/jarmo/secrets-cli.git
cd secrets-cli
make test
make install
```

`make install` builds the `secrets` executable into the Go binary directory. Use `INSTALL_DIR=/custom/bin` to change the installation directory.

Bump `VERSION` in `secrets.go`, commit the change, then run `make release`. The release script creates the GitHub tag automatically and fails if the tag already exists. Set `GITHUB_RELEASE_TOKEN` to a token with repository Contents write permission.

