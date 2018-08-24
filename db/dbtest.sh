#!/bin/sh

set -e

#if [ ! -f "build/env.sh" ]; then
    #echo "$0 must be run from the root of the repository."
    #exit 2
#fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
ethdir="$workspace/src/github.com/eosspark/eos-go"
if [ ! -L "$ethdir/db" ]; then
    mkdir -p "$ethdir"
    cd "$ethdir"
    ln -s ../../../../../../. db
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
#echo $GOPATH
export GOPATH

# Run the command inside the workspace.
cd "$ethdir/db"
PWD="$ethdir/db"

# Launch the arguments with the configured environment.
go test
cd $root
#echo $root
rm -rf ./eos.db
rm -rf ./build
#exec "$@"
