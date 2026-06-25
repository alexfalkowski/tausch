// Package io writes tausch-configured stdout/stderr payloads to an io.Writer.
//
// This package is internal to the tausch module. It bridges configuration values
// (which are stored as strings) to actual byte output written to a destination
// writer such as os.Stdout, os.Stderr, or a bytes.Buffer in tests.
//
// # Payload decoding
//
// The tausch YAML config stores `stdout` and `stderr` as strings in a
// `kind:data` format (for example `text:hello`, `file:/tmp/out.txt`,
// `base64:...`). This package decodes those payloads and copies the resulting
// stream to the supplied writer. Relative `file:` paths are resolved from the
// config file directory supplied by the caller.
//
// # Write semantics
//
// [Write] has a small but important contract:
//
//   - If data is the empty string, it performs no writes and returns (false, nil).
//   - If data is non-empty, it attempts to copy the decoded stream to w.
//     On success it returns (true, nil).
//   - If decoding or opening a file fails, it returns (false, err).
//   - If the underlying stream copy fails after output starts, it returns
//     (true, err).
//
// The returned boolean indicates whether output was attempted/emitted. The CLI
// orchestration uses this to choose between stdout and the fallback stderr path;
// configured exit status is handled outside this package.
//
// # Writer alias
//
// [Writer] is provided as a convenience alias to [io.Writer], primarily to avoid
// importing the standard library io package in callers that already use this
// internal package.
package io
