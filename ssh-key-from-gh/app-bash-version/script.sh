#!/bin/bash

organization="huetruong"
authorizedKeysFile="$HOME/.ssh/authorized_keys"

createAuthorizedKeysFile() {
    # Check if the file already exists
    if [ -f "$1" ]; then
        return
    fi

    # Create the directory if it doesn't exist
    mkdir -p "$(dirname "$1")"

    # Create the "authorized_keys" file
    touch "$1"
}

installPackage() {
    package="$1"

    # Check if the package is already installed
    if command -v "$package" >/dev/null 2>&1; then
        return
    fi

    echo "Installing $package..."
    if [ "$(uname)" == "Darwin" ]; then
        # macOS
        brew install "$package"
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        # Linux
        sudo apt-get update
        sudo apt-get install -y "$package"
    else
        echo "Error: Unsupported operating system"
        exit 1
    fi
}

getPublicKeys() {
    account="$1"
    url="https://api.github.com/users/$account/keys"

    response=$(curl -s "$url")
    status=$(echo "$response" | jq -r '. | length')

    if [ "$status" -ne 0 ]; then
        keys=$(echo "$response" | jq -r '.[].key')
        echo "$keys"
    else
        echo "Error: Request returned non-ok status"
        return 1
    fi
}

appendKeysToFile() {
    keys="$1"
    account="$2"
    filename="$3"

    for key in $keys; do
        comment="# Key ID: $key_id, User: $account"
        echo -e "$comment\n$key" >> "$filename"
    done
}

getOrganizationMembers() {
    organization="$1"
    url="https://api.github.com/orgs/$organization/members"

    response=$(curl -s "$url")
    status=$(echo "$response" | jq -r '. | length')

    if [ "$status" -ne 0 ]; then
        members=$(echo "$response" | jq -r '.[].login')
        echo "$members"
    else
        echo "Error: Request returned non-ok status"
        return 1
    fi
}

# Check and install required packages
installPackage "curl"
installPackage "jq"

# Create the "authorized_keys" file if it doesn't exist
createAuthorizedKeysFile "$authorizedKeysFile"

# Get the list of GitHub accounts from the organization
accounts=$(getOrganizationMembers "$organization")

for account in $accounts; do
    keys=$(getPublicKeys "$account")
    if [ $? -ne 0 ]; then
        echo "Error getting public keys for account $account"
        continue
    fi

    appendKeysToFile "$keys" "$account" "$authorizedKeysFile"
    if [ $? -ne 0 ]; then
        echo "Error writing public keys to file $authorizedKeysFile"
    else
        echo "Public keys for account $account successfully copied to file $authorizedKeysFile"
    fi
done
