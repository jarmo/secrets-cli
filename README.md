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
  secrets list [FILTER] [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets add NAME [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets edit ID [--alias=VAULT_ALIAS | --path=VAULT_PATH]
  secrets delete ID [--alias=VAULT_PATH | --path=VAULT_PATH]
  secrets change-password [--alias=VAULT_PATH | --path=VAULT_PATH]
  secrets initialize --path=VAULT_PATH --alias=VAULT_ALIAS

Options:
  --alias VAULT_ALIAS    Optional vault alias.
  --path VAULT_PATH      Optional vault path. Defaults to the path in configuration.
  -h --help              Show this screen.
  -v --version           Show version.
```

### Initializing Vault

Vault needs to be initialized if there is going to be a default vault. Otherwise specifying `--path` or `--alias` with any command is supported. Initializing vault just stores location and alias to your vault into a configuration file (supporting [XDG Base Directory standard](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)):

```
$ secrets initialize --path /home/user/.secrets.json --alias main
Vault successfully configured at /home/user/.config/secrets/config.json and is ready to store your secrets!
```

### Adding a New Secret

Add your first secret:

```
$ secrets add "my secret"
Enter vault password: [enter secure passphrase and remember it]
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
$ secrets list
Enter vault password: [your secure passphrase]

[299ed462-b171-4d67-ba21-264b221d9913]
my secret
my secret value
```

### Listing Specific Secrets

```
$ secrets list "secret"
Enter vault password: [your secure passphrase]

[299ed462-b171-4d67-ba21-264b221d9913]
my secret
my secret value
```

### Editing a Secret

```
$ secrets edit 299ed462-b171-4d67-ba21-264b221d9913
Enter vault password: [your secure passphrase]
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
$ secrets delete 299ed462-b171-4d67-ba21-264b221d9913
Enter vault password: 
Deleted: 
[299ed462-b171-4d67-ba21-264b221d9913]
different secret name
different secret value
yet another secret value line
```

## Using multiple vaults

Just append `--alias` after any command to operate against selected vault.
When `--alias` is not specified a first vault existing in configuration file will be used.

## But how do I sync vault between different devices?!

One way to sync would be to use any already existing syncing platforms like Dropbox, Microsoft OneDrive or Google Drive.
Since you can specify vault storage location then it is up to you how (or if even) you sync.

## Development

Retrieve dependencies, build and install binaries to `$GOPATH/bin/`

```
git clone https://github.com/jarmo/secrets-cli.git
cd secrets-cli
make
make install
```
