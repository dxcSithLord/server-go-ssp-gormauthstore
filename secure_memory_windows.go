//go:build windows

package gormauthstore

import (
	"reflect"
	"runtime"
	"syscall"
	"unsafe"

	ssp "github.com/sqrldev/server-go-ssp"
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

// ScrambleBytes overwrites a byte slice with pseudo-random data.
//
//go:noinline
func ScrambleBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	pattern := byte(0xAA)
	for i := range b {
		b[i] = pattern ^ byte(i&0xFF)
	}

	runtime.KeepAlive(b)
}

// WipeString attempts to securely clear a string by accessing its underlying bytes.
// See secure_memory.go for full documentation and warnings.
func WipeString(s *string) {
	if s == nil || *s == "" {
		return
	}

	sh := (*reflect.StringHeader)(unsafe.Pointer(s))
	if sh.Len == 0 || sh.Data == 0 {
		*s = ""
		return
	}

	sl := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}

	b := *(*[]byte)(unsafe.Pointer(&sl))
	WipeBytes(b)

	*s = ""
}

// ClearIdentity securely wipes all sensitive fields from a SqrlIdentity struct.
// See secure_memory.go for full documentation.
func ClearIdentity(identity *ssp.SqrlIdentity) {
	if identity == nil {
		return
	}

	WipeString(&identity.Idk)
	WipeString(&identity.Suk)
	WipeString(&identity.Vuk)
	WipeString(&identity.Pidk)
	WipeString(&identity.Rekeyed)

	identity.SQRLOnly = false
	identity.Hardlock = false
	identity.Disabled = false
	identity.Btn = 0

	runtime.KeepAlive(identity)
}

// SecureIdentityWrapper provides RAII-style automatic cleanup for SqrlIdentity.
type SecureIdentityWrapper struct {
	Identity *ssp.SqrlIdentity
	wiped    bool
}

// NewSecureIdentityWrapper creates a new wrapper around an existing identity.
func NewSecureIdentityWrapper(identity *ssp.SqrlIdentity) *SecureIdentityWrapper {
	return &SecureIdentityWrapper{
		Identity: identity,
		wiped:    false,
	}
}

// Destroy securely wipes the identity and marks the wrapper as invalid.
func (w *SecureIdentityWrapper) Destroy() {
	if w == nil || w.wiped {
		return
	}

	if w.Identity != nil {
		ClearIdentity(w.Identity)
		w.Identity = nil
	}
	w.wiped = true
}

// IsValid returns true if the wrapper still contains a valid identity.
func (w *SecureIdentityWrapper) IsValid() bool {
	return w != nil && !w.wiped && w.Identity != nil
}

// GetIdentity returns the wrapped identity if valid, otherwise returns nil.
func (w *SecureIdentityWrapper) GetIdentity() *ssp.SqrlIdentity {
	if !w.IsValid() {
		return nil
	}
	return w.Identity
}

// ValidateIdk performs basic validation on an Identity Key.
func ValidateIdk(idk string) error {
	if idk == "" {
		return ErrEmptyIdentityKey
	}

	if len(idk) > MaxIdkLength {
		return ErrIdentityKeyTooLong
	}

	for _, c := range idk {
		if !isValidIdkChar(c) {
			return ErrInvalidIdentityKeyFormat
		}
	}

	return nil
}

// isValidIdkChar checks if a character is valid for an Identity Key.
func isValidIdkChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '+' || c == '/' || c == '=' || c == '-' || c == '_' || c == '.'
}

// Constants for validation
const (
	MaxIdkLength = 256
)
