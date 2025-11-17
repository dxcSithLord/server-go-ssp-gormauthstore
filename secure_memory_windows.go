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

// WipeBytes securely zeros out a byte slice using Windows RtlSecureZeroMemory
// RtlSecureZeroMemory is a compiler-intrinsic that ensures memory is zeroed
// without being optimized away
func WipeBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// Call RtlSecureZeroMemory - it's void and never fails, so don't check return
	procSecureZero.Call(
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
	)

	// Prevent the slice from being garbage-collected before wiping completes
	runtime.KeepAlive(b)
}
