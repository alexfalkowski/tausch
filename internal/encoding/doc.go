// Package encoding decodes tausch "kind:data" payload strings into bytes.
//
// This package is internal to the tausch module. It is used to interpret the
// `stdout` and `stderr` fields of a configured command.
//
// # Format
//
// Values are encoded as a single string with a kind prefix and a colon:
//
//	kind:data
//
// The substring before the first ":" is the kind. Everything after the first ":"
// is the data portion (and may itself contain additional ":" characters).
//
// # Supported kinds
//
//   - text:<literal text>
//     Interprets data as UTF-8 (or arbitrary) text and returns its raw bytes.
//
//   - file:<path>
//     Interprets data as a filesystem path and returns the contents of the file
//     via os.ReadFile.
//
//   - base64:<base64-encoded bytes>
//     Interprets data as a standard base64 encoding (RFC 4648) and decodes it
//     using base64.StdEncoding.
//
// # Errors
//
// If the kind is not one of the supported values, [Decode] returns [ErrKindNotFound].
//
// Any errors encountered while reading a file or decoding base64 are returned
// directly from the underlying standard library functions.
package encoding
