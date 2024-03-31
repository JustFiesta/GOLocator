#!/usr/bin/env bash
# Setup script for GOLocator

# Functions start
# Cleanup function if installaion goes wrong
cleanup () {
    rm -rf GOLocator/
}


# Script start
# Create GOlocator folder and change directory to it
mkdir GOLocator
cd GOLocator

# Try to clone repo into GOLocator/
if [ $? -eq 0 ]; then
    echo "Cloning repository..."
    git clone https://github.com/JustFiesta/GOLocator

else # catch status for creating/changing dir to GOLocaor/
    echo "Failed to create and change to GOLocator folder!"
    exit 1
fi

# Try to build goloc binaries 
if [ $? -eq 0 ]; then
    echo "Building binaries..."
    go build -o goloc

else # catch status for git clone
    echo "Failed to clone repository!"
    exit 1
fi

# Try to install goloc binaries 
if [ $? -eq 0 ]; then
    echo "Installing binaries..."
    go install

else  # catch status for go build -o goloc
    echo "Failed to create build binary!"
    exit 1
fi

# Catch status of goloc insatllation
if [[ ! $? -eq 0 ]]; then
    echo "Failed to create build binary!"
    cleanup
    exit 1
fi