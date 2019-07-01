#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
eosgodir="$workspace/src/github.com/eosspark"
if [ ! -L "$eosgodir/eos-go" ]; then
    mkdir -p "$eosgodir"
    cd "$eosgodir"
    ln -s ../../../../../. eos-go
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$eosgodir/eos-go"
PWD="$eosgodir/eos-go"

# Launch the arguments with the configured environment.
exec "$@"

