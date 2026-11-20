package vaultfile

import (
	"crypto/ecdh"
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/input"
	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/user"
	"github.com/jarmo/secrets/secret"
	"github.com/jarmo/secrets/storage"
)

type Vault struct {
	Secrets    []secret.Secret
	path       string
	password   []byte
	privateKey *ecdh.PrivateKey
}

func Open(vaultPath, username, usersDirectory string) Vault {
	if username == "" {
		secrets, path, password := Read(vaultPath)
		return Vault{Secrets: secrets, path: path, password: password}
	}
	password := input.AskPassword("Enter user password: ")
	defer clearPassword(password)
	vault, err := OpenWithPassword(vaultPath, username, usersDirectory, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return vault
}

func OpenWithPassword(vaultPath, username, usersDirectory string, password []byte) (Vault, error) {
	if username == "" {
		secrets, err := storage.Read(vaultPath, password)
		return Vault{Secrets: secrets, path: vaultPath, password: password}, err
	}
	if usersDirectory == "" {
		usersDirectory = path.UsersDirectory()
	}
	_, privateKey, err := user.Authenticate(usersDirectory, username, password)
	if err != nil {
		return Vault{}, err
	}
	secrets, err := storage.ReadForUser(vaultPath, privateKey)
	if err != nil {
		return Vault{}, err
	}
	return Vault{Secrets: secrets, path: vaultPath, privateKey: privateKey}, nil
}

func (vault Vault) Write(secrets []secret.Secret) error {
	if vault.privateKey != nil {
		return storage.WriteForUser(vault.path, vault.privateKey, secrets)
	}
	storage.Write(vault.path, vault.password, secrets)
	return nil
}

func (vault Vault) PrivateKey() *ecdh.PrivateKey {
	return vault.privateKey
}

func clearPassword(password []byte) {
	for i := range password {
		password[i] = 0
	}
}
