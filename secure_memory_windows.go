//go:build windows

package gormauthstore

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	procSecureZero = kernel32.NewProc("RtlSecureZeroMemory")
)

// WipeBytes securely overwrites a byte slice with zeros using Windows
// RtlSecureZeroMemory, which is guaranteed not to be optimized away by the compiler.
//
//go:noinline
func WipeBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// Try to use Windows secure memory clearing
	ret, _, _ := procSecureZero.Call(
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
	)

	// Fallback if system call fails
	if ret == 0 {
		for i := range b {
			b[i] = 0
		}
	}

	runtime.KeepAlive(b)
}
