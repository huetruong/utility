#!/bin/bash
# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go..."

    # Download the Go binary archive
    wget https://golang.org/dl/go1.17.5.linux-amd64.tar.gz

    # Extract the archive
    sudo tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz

    # Add Go binaries to the PATH environment variable
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile.d/go.sh

    # Reload the profile
    source /etc/profile.d/go.sh

    rm go1.17.5.linux-amd64.tar.gz
    
    echo "Go has been installed."
fi


# Clone the repository
echo "Cloning the repository..."
git clone https://github.com/JustYAMLGuys/utility.git /tmp/app
cd /tmp/app/utility/ssh-key-from-gh/app || exit 1

# Build and run the Go program
echo "Building and running the program..."
sudo -u ubuntu /usr/local/go/bin/go run main.go

# Cleanup: Remove the cloned repository
echo "Cleaning up..."
cd ../..
rm -rf /tmp/app

echo "Script completed successfully."