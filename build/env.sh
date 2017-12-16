#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
zrmdir="$workspace/src/github.com/apolo-technologies"
if [ ! -L "$zrmdir/zerium" ]; then
    mkdir -p "$zrmdir"
    cd "$zrmdir"
    ln -s ../../../../../. zerium
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$zrmdir/zerium"
PWD="$zrmdir/zerium"

# Launch the arguments with the configured environment.
exec "$@"
