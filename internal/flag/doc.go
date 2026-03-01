// Package flag implements tausch CLI argument parsing and config path resolution.
//
// This package is internal to the tausch module. It is responsible for
// interpreting the arguments passed to the `tausch` binary, extracting the
// configuration file location, and deriving the command name used to look up an
// entry in the YAML configuration.
//
// # Parsing model
//
// The tausch CLI is designed to be invoked in the following shape:
//
//	tausch -config path/to/config.yml -- <your command tokens...>
//
// This package uses the standard library's flag parsing to consume recognized
// flags. Any remaining arguments (those after flag parsing, typically the tokens
// after `--`) are treated as the command to stub.
//
// # Command name derivation
//
// The command "name" is derived by joining the remaining arguments with a single
// space and trimming leading/trailing whitespace.
//
// For example, for:
//
//	tausch -- go version
//
// the derived name is:
//
//	"go version"
//
// This derived name must match the `name` field in the YAML config exactly.
// Because matching is string-based, differences in spacing or quoting will cause
// lookups to fail.
//
// # Config path resolution
//
// The config file path is resolved in the following precedence order:
//
//  1. The `-config` flag, if provided.
//  2. The `TAUSCH_CONFIG` environment variable, if set.
//  3. The default path under the user config directory:
//     `$HOME/.config/tausch/config.yml` (or the platform equivalent returned by
//     os.UserConfigDir).
//
// The resolver only selects a path; it does not verify the file exists or is
// readable. Those checks occur when the caller attempts to open/decode the file.
package flag
