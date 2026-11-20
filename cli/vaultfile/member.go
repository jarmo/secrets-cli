package vaultfile

import (
	"crypto/ecdh"
	"fmt"
	"os"

	"github.com/jarmo/secrets-cli/v6/input"
	"github.com/jarmo/secrets-cli/v6/storage/path"
	"github.com/jarmo/secrets/member"
	"github.com/jarmo/secrets/secret"
	"github.com/jarmo/secrets/storage"
)

type Vault struct {
	Secrets    []secret.Secret
	path       string
	password   []byte
	privateKey *ecdh.PrivateKey
}

func Open(vaultPath, username, membersDirectory string) Vault {
	if username == "" {
		secrets, path, password := Read(vaultPath)
		return Vault{Secrets: secrets, path: path, password: password}
	}
	password := input.AskPassword("Enter member password: ")
	defer clearPassword(password)
	vault, err := OpenWithPassword(vaultPath, username, membersDirectory, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return vault
}

func OpenWithPassword(vaultPath, username, membersDirectory string, password []byte) (Vault, error) {
	if username == "" {
		secrets, err := storage.Read(vaultPath, password)
		return Vault{Secrets: secrets, path: vaultPath, password: password}, err
	}
	if membersDirectory == "" {
		membersDirectory = path.MembersDirectory()
	}
	_, privateKey, err := member.Authenticate(membersDirectory, username, password)
	if err != nil {
		return Vault{}, err
	}
	secrets, err := storage.ReadForMember(vaultPath, privateKey)
	if err != nil {
		return Vault{}, err
	}
	return Vault{Secrets: secrets, path: vaultPath, privateKey: privateKey}, nil
}

func (vault Vault) Write(secrets []secret.Secret) error {
	if vault.privateKey != nil {
		return storage.WriteForMember(vault.path, vault.privateKey, secrets)
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
