include bin/build/make/help.mak
include bin/build/make/go.mak
include bin/build/make/git.mak

# Build the cli.
build:
	@go build
