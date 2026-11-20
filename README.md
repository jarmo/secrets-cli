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
  secrets list [FILTER] [--member=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets add NAME [--member=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets edit ID [--member=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets delete ID [--member=USERNAME] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets change-password [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets initialize --path=VAULT_PATH --alias=VAULT_ALIAS --member=USERNAME [--members-dir=MEMBERS_DIR]
  secrets member-create MEMBER [--members-dir=MEMBERS_DIR]
  secrets member-activate MEMBER [--members-dir=MEMBERS_DIR]
  secrets member-add MEMBER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--members-dir=MEMBERS_DIR] [--member=USERNAME]
  secrets member-list [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--members-dir=MEMBERS_DIR] [--member=USERNAME]
  secrets member-remove MEMBER [--alias=VAULT_ALIAS | --path=VAULT_PATH] [--members-dir=MEMBERS_DIR] [--member=USERNAME]

Options:
  --alias VAULT_ALIAS    Vault alias.
  --path VAULT_PATH      Vault path.
  --members-dir MEMBERS_DIR    Member storage directory.
  --member USERNAME    Authenticate using this member account.
  -h --help              Show this screen.
  -v --version           Show version.
```

### Initializing Vault

After creating and activating a member account (see **Member setup** below), create a vault and save its location and alias to a configuration file (supporting [XDG Base Directory standard](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)):

```
$ secrets initialize --path /home/user/.secrets.json --alias main --member alice
Enter member password: [Alice's password]
Vault successfully created and configured at /home/user/.config/secrets/config.json!
```

### Adding a New Secret

Add your first secret:

```
$ secrets add "my secret" --member alice
Enter member password: [Alice's password]
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
$ secrets list --member alice
Enter member password: [Alice's password]

[299ed462-b171-4d67-ba21-264b221d9913]
my secret
my secret value
```

### Listing Specific Secrets

```
$ secrets list "secret" --member alice
Enter member password: [Alice's password]

[299ed462-b171-4d67-ba21-264b221d9913]
my secret
my secret value
```

### Editing a Secret

```
$ secrets edit 299ed462-b171-4d67-ba21-264b221d9913 --member alice
Enter member password: [Alice's password]
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
$ secrets delete 299ed462-b171-4d67-ba21-264b221d9913 --member alice
Enter member password: 
Deleted: 
[299ed462-b171-4d67-ba21-264b221d9913]
different secret name
different secret value
yet another secret value line
```

### Member setup

Create a member and share the generated setup code securely:

```
$ secrets member-create bob
Setup code for bob (expires ...):
[setup code]
```

Bob activates the member account by choosing a password:

```
$ secrets member-activate bob
Enter setup code: [setup code]
Enter new member password: [password]
Enter new member password again: [password]
Member bob activated!
```

Setup codes expire after 48 hours. Both commands accept `--members-dir` to use the same member directory.

Once activated, a member can be granted access to a vault:

```
$ secrets member-add bob --alias production --member alice
Enter member password: [Alice's password]
Member bob added to vault!
```

`member-add` also accepts `--path` and `--members-dir`.

Once granted access, a member can use their own credentials with `--member`:

```
$ secrets list --alias production --member bob
Enter member password: [Bob's password]
```

The same option works with `add`, `edit`, `delete`, `member-add`, `member-list` and `member-remove`. Use `--members-dir` if the accounts are stored outside the default directory.

List members or revoke a member's access to a vault:

```
$ secrets member-list --alias production --member alice
alice
bob

$ secrets member-remove bob --alias production --member alice
Enter member password: [Alice's password]
Member bob removed from vault!
```

Removing a member rotates the vault encryption key while retaining access for other members. Both commands accept `--path` and `--members-dir`. The member's account is not deleted; it can be added to a vault again. Old backups can still be opened with keys that were valid at the time.

To revoke a member from all configured vaults and delete the account:

```
$ secrets member-destroy bob --member alice
Enter member password: [Alice's password]
Member bob destroyed!
```

The acting member must have access to every configured vault. The command checks all vaults before making changes and does not delete the account unless all revocations succeed. All vaults containing the member must be in the CLI configuration. Backups retain historical access.

## Using multiple vaults

Just append `--alias` after any command to operate against selected vault.
When `--alias` is not specified a first vault existing in configuration file will be used.

## But how do I sync vault between different devices?!

One way to sync would be to use any already existing syncing platforms like Dropbox, Microsoft OneDrive or Google Drive.
Since you can specify vault storage location then it is up to you how (or if even) you sync.

## Legacy vaults

Starting with version 7.x, vaults use a new format that supports sharing between users. Legacy vaults are automatically converted to the new format when written. Support for the legacy format will be removed in the next major version.

To share a legacy vault, run `secrets member-add alice --alias production` and enter the original vault password. This grants Alice access and retires the vault password. For vaults already shared with members, running `member-add` for an existing member also retires any remaining vault-password access.

## Development

Retrieve dependencies, build and install binaries to `$GOPATH/bin/`

```
git clone https://github.com/jarmo/secrets-cli.git
cd secrets-cli
make
make install
```
