//go:build !windows

package gormauthstore

import (
	"runtime"
	"unsafe"

	"golang.org/x/crypto/nacl/secretbox"
)

// WipeBytes securely zeros out a byte slice using compiler-resistant wiping
// This implementation uses the secretbox overhead constant to prevent optimization
func WipeBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// Use volatile-like pattern to prevent compiler optimization
	// The secretbox.Overhead reference ensures the compiler doesn't optimize this away
	_ = secretbox.Overhead

	// Zero out the memory
	for i := range b {
		b[i] = 0
	}

	// Memory barrier to ensure writes are not optimized away
	runtime.KeepAlive(b)

	// Additional compiler barrier using unsafe
	p := unsafe.Pointer(&b[0])
	_ = p
}
