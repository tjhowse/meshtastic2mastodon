#!/bin/bash -e

# https://protobuf.dev/getting-started/gotutorial/
# Install protoc from https://github.com/protocolbuffers/protobuf/releases
# You've gotta download and extract the zip, then put it on your PATH somewhere. Also there are some
# bundled "stdlib" proto files that need to be in a "include" directory next to the "bin" directory.
# Run go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Get a temp directory
TMPDIR=$(mktemp -d)
echo "Using temp dir: $TMPDIR"

# Set up a trap to clean up the temp directory on exit
trap "rm -rf $TMPDIR" EXIT

# Clone the repo
git clone https://github.com/meshtastic/protobufs.git $TMPDIR

protoc -I=$TMPDIR --go_out=. $TMPDIR/nanopb.proto $TMPDIR/meshtastic/*.proto

rm -rf ./generated

mv github.com/meshtastic/go/generated .

rm -rf github.com