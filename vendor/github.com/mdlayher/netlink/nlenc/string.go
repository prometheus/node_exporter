package nlenc

import "bytes"

// Bytes returns a null-terminated byte slice with the contents of s.
func Bytes(s string) []byte {
	return append([]byte(s), 0x00)
}

// String returns a string with the contents of b from a null-terminated
// byte slice.
func String(b []byte) string {
	return string(bytes.TrimSuffix(b, []byte{0x00}))
}
