// Package io writes tausch-configured stdout/stderr payloads to an io.Writer.
//
// This package is internal to the tausch module. It bridges configuration values
// (which are stored as strings) to actual byte output written to a destination
// writer such as os.Stdout, os.Stderr, or a bytes.Buffer in tests.
//
// # Relationship to internal/encoding
//
// The tausch YAML config stores `stdout` and `stderr` as strings in a
// `kind:data` format (for example `text:hello`, `file:/tmp/out.txt`,
// `base64:...`). This package delegates decoding of that representation to
// internal/encoding and then writes the resulting bytes to the supplied writer.
//
// # Write semantics
//
// [Write] has a small but important contract:
//
//   - If data is the empty string, it performs no writes and returns (false, nil).
//   - If data is non-empty, it attempts to decode it and write it to w.
//     On success it returns (true, nil).
//   - If decoding fails or the underlying write fails, it returns (false, err).
//
// The returned boolean indicates whether output was attempted/emitted. The CLI
// orchestration uses this to decide whether to treat a command as a "success"
// path (stdout present) or "error" path (fallback to stderr + non-zero exit).
//
// # Writer alias
//
// [Writer] is provided as a convenience alias to [io.Writer], primarily to avoid
// importing the standard library io package in callers that already use this
// internal package.
package io
