include bin/build/make/help.mak
include bin/build/make/go.mak
include bin/build/make/git.mak

# Build the cli.
build:
	@go build

# Run all repository Go benchmarks.
benchmarks: benchmark-cli benchmark-config benchmark-flag benchmark-io

# Benchmark the end-to-end CLI stub path.
benchmark-cli:
	@$(MAKE) benchmark

# Benchmark config loading and command lookup.
benchmark-config:
	@$(MAKE) benchmark package=internal/config

# Benchmark CLI flag parsing and command-name derivation.
benchmark-flag:
	@$(MAKE) benchmark package=internal/flag

# Benchmark configured output decoding and writing.
benchmark-io:
	@$(MAKE) benchmark package=internal/io

# Run bounded smoke fuzzing for all repository fuzz targets.
fuzzes: fuzz-cmd fuzz-config fuzz-flag fuzz-io fuzz-exec

# Fuzz the end-to-end CLI stdout path.
fuzz-cmd:
	@$(MAKE) fuzz package=internal/cmd name=FuzzRunWritesConfiguredStdout fuzztime=1s

# Fuzz config validation and command lookup.
fuzz-config:
	@$(MAKE) fuzz package=internal/config name=FuzzConfigValidate fuzztime=1s
	@$(MAKE) fuzz package=internal/config name=FuzzConfigGetCommand fuzztime=1s

# Fuzz CLI argument parsing and command-name derivation.
fuzz-flag:
	@$(MAKE) fuzz package=internal/flag name=FuzzParseCommandName fuzztime=1s

# Fuzz output decoding and writing.
fuzz-io:
	@$(MAKE) fuzz package=internal/io name=FuzzWriteText fuzztime=1s
	@$(MAKE) fuzz package=internal/io name=FuzzWriteBase64 fuzztime=1s
	@$(MAKE) fuzz package=internal/io name=FuzzWriteKindNotFound fuzztime=1s

# Fuzz the public exec wrapper delimiter contract.
fuzz-exec:
	@$(MAKE) fuzz package=exec name=FuzzCommandPrefixesDelimiter fuzztime=1s
