//go:build !windows

package gormauthstore

import (
	"runtime"
)

// WipeBytes securely overwrites a byte slice with zeros.
// This function uses compiler directives to prevent dead store elimination
// which could otherwise optimise away the memory clearing operation.
//
// Note: This provides best-effort clearing but Go's garbage collector may still
// have copies of data in memory. For maximum security, consider using memguard
// library with locked memory pages.
//
//go:noinline
func WipeBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// Use a simple loop to zero out the bytes
	for i := range b {
		b[i] = 0
	}

	// Force a reference to prevent dead store elimination
	// This ensures the compiler doesn't optimise away our zeroing
	runtime.KeepAlive(b)
}
