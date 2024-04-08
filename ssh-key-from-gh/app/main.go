package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	// GitHub organization name
	organization := "huemattic"

	// Path to the "authorized_keys" file
	authorizedKeysFile := filepath.Join(getUserHomeDir(), ".ssh", "authorized_keys")

	// Create the "authorized_keys" file if it doesn't exist
	err := createAuthorizedKeysFile(authorizedKeysFile)
	if err != nil {
		log.Printf("Error creating authorized_keys file: %s\n", err.Error())
		return
	}

	// Get the list of GitHub accounts from the organization
	accounts, err := getOrganizationMembers(organization)
	if err != nil {
		log.Printf("Error getting organization members: %s\n", err.Error())
		return
	}

	// Get the public keys for each account and append them to the "authorized_keys" file
	for _, account := range accounts {
		keys, err := getPublicKeys(account)
		if err != nil {
			log.Printf("Error getting public keys for account %s: %s\n", account, err.Error())
			continue
		}

		err = appendKeysToFile(keys, account, authorizedKeysFile)
		if err != nil {
			log.Printf("Error writing public keys to file %s: %s\n", authorizedKeysFile, err.Error())
		} else {
			log.Printf("Public keys for account %s successfully copied to file %s\n", account, authorizedKeysFile)
		}
	}
}

// Data structure to represent a public key
type PublicKey struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

// Get public keys from a GitHub account
func getPublicKeys(account string) ([]PublicKey, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/keys", account)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request returned non-ok status: %s", resp.Status)
	}

	var keys []PublicKey
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err.Error())
	}

	return keys, nil
}

// Append public keys to the "authorized_keys" file
func appendKeysToFile(keys []PublicKey, account, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, key := range keys {
		comment := fmt.Sprintf("# Key ID: %d, User: %s\n", key.ID, account)
		_, err := file.WriteString(comment + key.Key + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Get the list of members from a GitHub organization
func getOrganizationMembers(organization string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members", organization)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request returned non-ok status: %s", resp.Status)
	}

	var members []struct {
		Login string `json:"login"`
	}

	err = json.NewDecoder(resp.Body).Decode(&members)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err.Error())
	}

	var accounts []string
	for _, member := range members {
		accounts = append(accounts, member.Login)
	}

	return accounts, nil
}

// Create the "authorized_keys" file if it doesn't exist
func createAuthorizedKeysFile(filename string) error {
	// Check if the file already exists
	_, err := os.Stat(filename)
	if err == nil {
		return nil // File already exists, no need to create
	}

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(filename), 0700)
	if err != nil {
		return err
	}

	// Create the "authorized_keys" file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// Get the current user's home directory
func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir
}
